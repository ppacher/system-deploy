package deploy

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ppacher/system-conf/conf"
	"github.com/ppacher/system-deploy/pkg/condition"
)

// ErrInvalidTaskSection is returned if a section is invalid.
var ErrInvalidTaskSection = errors.New("invalid task section")

// Task defines a deploy task.
type Task struct {
	file *conf.File

	// FileName is the name of the file that describes this task.
	FileName string

	// Directory holds the directory of the task.
	Directory string

	// Description is the tasks description.
	Description string

	// StartMasked is set to true if this task is disabled (masked)
	// by default.
	StartMasked bool

	// Disabled can be set to true to disable a task permanently.
	Disabled bool

	// EnvironmentFiles holds a list of environment files.
	EnvironmentFiles []string

	// Sections holds the tasks sections.
	Sections []conf.Section

	// Environment holds the parsed environment.
	Environment []string

	// Conditions is a list of conditions that must match.
	Conditions []condition.Instance
}

// DecodeFile is like Decode but reads the task from
// filePath.
func DecodeFile(filePath string) (*Task, error) {
	tsk, err := decodeFile(filePath)
	return tsk, err
}

func decodeFile(filePath string) (*Task, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return decode(filePath, f)
}

// Decode decodes a deploy task from r and uses the basename
// of name as the task's name.
func Decode(filePath string, r io.Reader) (*Task, error) {
	tsk, err := decode(filePath, r)
	return tsk, err
}

func decode(filePath string, r io.Reader) (*Task, error) {
	file, err := conf.Deserialize(filePath, r)
	if err != nil {
		return nil, err
	}

	task := &Task{
		file:      file,
		FileName:  filepath.Base(filePath),
		Directory: filepath.Dir(filePath),
		Sections:  file.Sections,
	}

	if err := applyMetaData(task); err != nil {
		return nil, err
	}

	if len(task.Sections) == 0 {
		// we must return task and metaSection here because
		// ErrNoSections is ignored when loading drop-ins.
		return task, conf.ErrNoSections
	}

	return task, nil
}

func applyMetaData(task *Task) error {
	for idx, sec := range task.file.Sections {
		if strings.ToLower(sec.Name) == "task" {

			if err := decodeMetaData(sec, task); err != nil {
				return ErrInvalidTaskSection
			}

			task.Sections = append(task.file.Sections[:idx], task.file.Sections[idx+1:]...)

			break
		}
	}

	return nil
}

func decodeMetaData(section conf.Section, task *Task) error {
	if strings.ToLower(section.Name) != "task" {
		return errors.New("invalid section name")
	}

	specs := make([]conf.OptionSpec, len(taskOptions))
	for idx, spec := range taskOptions {
		vals := section.Options.GetStringSlice(spec.Name)
		if len(vals) > 0 && spec.set != nil {
			if err := spec.set(section.Options, task); err != nil {
				return err
			}
		}

		specs[idx] = spec.OptionSpec
	}

	return conf.ValidateOptions(section.Options, specs)
}

// Clone creates a deep copy of t.
func (tsk *Task) Clone() *Task {
	n := &Task{
		FileName:    tsk.FileName,
		Directory:   tsk.Directory,
		Description: tsk.Description,
		StartMasked: tsk.StartMasked,
		Disabled:    tsk.Disabled,
	}

	if tsk.EnvironmentFiles != nil {
		n.EnvironmentFiles = make([]string, len(tsk.EnvironmentFiles))
		copy(n.EnvironmentFiles, tsk.EnvironmentFiles)
	}

	if tsk.Environment != nil {
		n.Environment = make([]string, len(tsk.Environment))
		copy(n.Environment, tsk.Environment)
	}

	if tsk.Conditions != nil {
		n.Conditions = make([]condition.Instance, len(tsk.Conditions))
		copy(n.Conditions, tsk.Conditions)
	}

	if len(tsk.Sections) > 0 {
		n.Sections = make([]conf.Section, len(tsk.Sections))
		for idx, s := range tsk.Sections {
			n.Sections[idx] = conf.Section{
				Name:    s.Name,
				Options: make(conf.Options, len(s.Options)),
			}

			for optIdx, opt := range s.Options {
				n.Sections[idx].Options[optIdx] = conf.Option{
					Name:  opt.Name,
					Value: opt.Value,
				}
			}
		}
	} else if tsk.Sections != nil {
		// make sure we also have an empty slice
		n.Sections = make([]conf.Section, 0)
	}

	return n
}
