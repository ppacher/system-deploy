package copy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	copyDir "github.com/otiai10/copy"

	"github.com/ppacher/system-deploy/pkg/actions"
	"github.com/ppacher/system-deploy/pkg/change"
	"github.com/ppacher/system-deploy/pkg/deploy"
	"github.com/ppacher/system-deploy/pkg/unit"
	"github.com/ppacher/system-deploy/pkg/utils"
)

func init() {
	actions.MustRegister(actions.Plugin{
		Name:        "Copy",
		Description: "Copy files and folder and script updates",
		Setup:       setupAction,
		Example:     example,
		Author:      "Patrick Pacher <patrick.pacher@gmail.com>",
		Website:     "https://github.com/ppacher/system-deploy",
		Help: []deploy.Section{
			{
				Title: "Change Detection",
				Description: "" +
					"The [Copy] action uses a Murmur3 hash to check whether or not a destination file needs to be updated. " +
					"In any case, [Copy] ensures the destination files mode bit either match the value of FileMode= or the mode bits of the source file. " +
					"See FileMode= for more information.",
			},
			{
				Title: "Bugs",
				Description: "Note that " + color.New(color.Bold).Sprint("FileMode") + " does not work when copying a directory recursively. " +
					"Also, directories are always copied without checking if an update is required." +
					"This will be fixed in a later release.",
			},
		},
		Options: []deploy.OptionSpec{
			{
				Name:        "Source",
				Required:    true,
				Description: "The source file to copy to Destination.",
				Type:        deploy.StringType,
			},
			{
				Name:        "Destination",
				Required:    true,
				Description: "The destination path where Source should be copied to.",
				Type:        deploy.StringType,
			},
			{
				Name:        "CreateDirectories",
				Description: "If set to true, missing directories in Destination will be created.",
				Type:        deploy.BoolType,
				Default:     "no",
			},
			{
				Name: "FileMode",
				Description: "" +
					"The mode bits (before umask) to use for the destination file. If unset the source files" +
					"mode bits will be used. The destination files mode will be changed to match FileMode= " +
					"even if the content is already correct." +
					"Note that Mode is ingnored when copying directories.",
				Type:    deploy.IntType,
				Default: "",
			},
			{
				Name:        "DirectoryMode",
				Description: "When creating Destination path (CreateDirectories=yes) the mode bits (before umask) for that directories.",
				Type:        deploy.IntType,
				Default:     "0755",
			},
		},
	})
}

func setupAction(task deploy.Task, sec unit.Section) (actions.Action, error) {
	a := &action{
		taskDir: task.Directory,
		opts:    sec.Options,
	}

	return a, nil
}

func (a *action) Prepare(graph actions.ExecGraph) error {
	{
		source, err := a.opts.GetString("Source")
		if err != nil {
			return err
		}
		a.source = filepath.Clean(source)

		if !filepath.IsAbs(a.source) {
			a.source = filepath.Clean(filepath.Join(a.taskDir, a.source))
		}
	}

	// get the destination path but don't clean it yet because
	// we need to know if the user specified a trailing path
	// separator.
	destination, err := a.opts.GetString("Destination")
	if err != nil {
		return err
	}

	{
		a.createPath, err = a.opts.GetBool("CreateDirectories")
		if err != nil && !unit.IsNotSet(err) {
			return err
		}
	}

	{
		fileMode, err := a.opts.GetInt("FileMode")
		if err != nil {
			if !unit.IsNotSet(err) {
				return fmt.Errorf("invalid value for FileMode: %w", err)
			}
			fileMode = 0700
		} else {
			if fileMode > 0777 {
				return fmt.Errorf("invalid value for FileMode: %o", fileMode)
			}
		}
		a.fileMode = os.FileMode(fileMode)
	}

	{
		dirMode, err := a.opts.GetInt("DirectoryMode")
		if err != nil {
			if !unit.IsNotSet(err) {
				return fmt.Errorf("invalid value for DirectoryMode: %w", err)
			}
			dirMode = 0755
		} else {
			if dirMode > 0777 {
				return fmt.Errorf("invalid value for DirectoryMode: %o", dirMode)
			}
		}
		a.dirMode = os.FileMode(dirMode)
	}

	{
		fi, err := os.Stat(a.source)
		if err != nil {
			return fmt.Errorf("source: %w", err)
		}
		a.sourceIsDir = fi.IsDir()
	}

	{
		if strings.HasSuffix(destination, string(filepath.Separator)) {
			a.destDir = destination
			a.destName = filepath.Base(a.source)
		} else {
			a.destDir = filepath.Dir(destination)
			a.destName = filepath.Base(destination)
		}

		if err := checkDirectory(a.destDir, a.createPath); err != nil {
			return err
		}
	}

	return nil
}

type action struct {
	actions.Base

	taskDir string
	opts    unit.Options
	log     actions.Logger

	source      string
	sourceIsDir bool
	fileMode    os.FileMode
	dirMode     os.FileMode
	destDir     string
	destName    string
	createPath  bool

	runPost bool
}

func (a *action) Name() string {
	return "Copy " + a.source + " to " + filepath.Join(a.destDir, a.destName)
}

func (a *action) Run(ctx context.Context) (bool, error) {
	var changed bool

	if a.createPath {
		if err := os.MkdirAll(a.destDir, a.dirMode); err != nil {
			return false, fmt.Errorf("failed to create destination %q: %w", a.destDir, err)
		}
	}

	dest := filepath.Join(a.destDir, a.destName)
	if a.sourceIsDir {
		// TODO(ppacher): allow specifying symlink actions.
		if err := copyDir.Copy(a.source, dest, copyDir.DefaultOptions); err != nil {
			return false, fmt.Errorf("failed to copy directory: %w", err)
		}
		changed = true // change detection is not yet supported when copying directories.
	} else {
		var err error
		changed, err = a.copyRegularFile()
		if err != nil {
			return false, err
		}
	}

	a.runPost = changed

	return changed, nil
}

func (a *action) getModeForFile() (os.FileMode, error) {
	fileMode := a.fileMode
	if fileMode == 0 {
		info, err := os.Lstat(a.source)
		if err != nil {
			return 0, fmt.Errorf("failed to stat source %q: %s", a.source, err)
		}

		fileMode = info.Mode()
	}
	return fileMode, nil
}

func (a *action) copyRegularFile() (bool, error) {
	dest := filepath.Join(a.destDir, a.destName)

	// find out which file mode we need to use, that is, either the one
	// specified via FileMode= or the mode bits of the source file.
	fileMode, err := a.getModeForFile()
	if err != nil {
		return false, err
	}

	// check if we actually need to update dest.
	updateRequired, err := change.FileUpdateNeeded(a.source, dest)
	if err != nil {
		return false, fmt.Errorf("failed to check for required file update: %w", err)
	}
	if !updateRequired {
		// file already exists and has the expected content, make sure
		// we have the correct file mode and we are done.
		return change.EnsureFileMode(dest, fileMode)
	}

	// finally replace/create dest from a.source and apply the correct
	// file mode. If dest exists it will be overwritten.
	if err := utils.CopyAtomicMode(a.source, dest, fileMode); err != nil {
		return false, err
	}

	return true, nil
}

func checkDirectory(path string, ignoreMissing bool) error {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if !ignoreMissing {
				return fmt.Errorf("path does not exist: %q", path)
			}

			return nil
		}
		return fmt.Errorf("failed to stat %q: %w", path, err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("not a directory: %q", path)
	}
	return nil
}

const example = `[Task]
Description= Copy file foo to /server/custom/bin

[Copy]
Source=./assets/foo
Destination=/server/custom/bin
CreateDirectories=yes
FileMode=0600`
