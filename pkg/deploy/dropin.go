package deploy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ppacher/system-deploy/pkg/unit"
)

// DropInExt is the file extension for drop-in files.
const DropInExt = ".conf"

// DropIn is a drop-in file for a given system-deploy task.
type DropIn struct {
	File string
	Task *Task
}

// readDir is used to read the contents of a directory and return
// a slice of os.FileInfo for each directory entry. It's here for
// unit-testing purposes and nomally points to ioutil.ReadDir.
var readDir func(path string) ([]os.FileInfo, error) = ioutil.ReadDir

// ApplyDropIns applies all dropins on t. DropIns can only be applied
// to tasks with unique section names. That is, if a task specifies
// the same action multiple times (like multiple [Copy] sections),
// drop-ins cannot be applied to that task.
func ApplyDropIns(t *Task, dropins []*DropIn, specs map[string]map[string]OptionSpec) (*Task, error) {
	copy := t.Clone()

	slm := make(map[string]*unit.Section)

	for _, sec := range copy.Sections {
		sn := strings.ToLower(sec.Name)
		if _, ok := slm[sn]; ok {
			// that section is defined multiple times
			// so instead of setting it we nil it.
			slm[sn] = nil
			continue
		}

		slm[sn] = &sec
	}

	for _, d := range dropins {
		if d.Task.StartMasked != nil {
			copy.StartMasked = &(*d.Task.StartMasked)
		}
		if d.Task.Disabled != nil {
			copy.Disabled = &(*d.Task.Disabled)
		}
		if d.Task.Description != nil {
			copy.Description = &(*d.Task.Description)
		}

		for _, dropInSec := range d.Task.Sections {
			sn := strings.ToLower(dropInSec.Name)

			s, ok := slm[sn]
			if !ok {
				return nil, ErrDropInSectionNotExists
			}

			sectionSpec, ok := specs[sn]
			if s == nil || !ok {
				return nil, ErrDropInSectionNotAllowed
			}

			// build a lookup map for the option values in this
			// drop-in section
			olm := make(map[string][]unit.Option)
			for _, opt := range dropInSec.Options {
				on := strings.ToLower(opt.Name)
				olm[on] = append(olm[on], opt)
			}

			// update each option, one after the other
			for optName, opts := range olm {
				optSpec, ok := sectionSpec[optName]
				if !ok {
					return nil, ErrOptionNotExists
				}

				// if the first value is empty it means we should
				// remove all current values in a slice type.
				// If it's not a slice type we are going to overwrite the existing
				// value so we can also remove it.
				if !optSpec.Type.IsSliceType() || opts[0].Value == "" {
					var newOpts unit.Options
					for _, opt := range s.Options {
						if strings.ToLower(opt.Name) != optName {
							newOpts = append(newOpts, opt)
						}
					}
					s.Options = newOpts

					if optSpec.Type.IsSliceType() {
						opts = opts[1:]
					}
				}

				// add the new values to the list
				s.Options = append(s.Options, opts...)
			}
		}
	}

	// rebuild the section slice.
	for idx, sec := range copy.Sections {
		copy.Sections[idx] = *slm[strings.ToLower(sec.Name)]
	}

	return copy, nil
}

// LoadDropIns loads all drop-in files for unitName. See SearchDropInFiles
// and DropInSearchPaths for more information on the searchPath.
func LoadDropIns(unitName string, searchPath []string) ([]*DropIn, error) {
	files, err := SearchDropinFiles(unitName, searchPath)
	if err != nil {
		return nil, err
	}

	dropins := make([]*DropIn, len(files))
	for idx, filePath := range files {
		t, err := DecodeFile(filePath)
		if err != nil {
			// don't ignore ErrNotExist here because
			// it existed just a few seconds ago!
			return nil, err
		}

		// Fix the filename to use unitName and
		// clear out the directory.
		t.FileName = unitName
		t.Directory = ""

		dropins[idx] = &DropIn{
			File: filePath,
			Task: t,
		}
	}

	return dropins, nil
}

// SearchDropinFiles searches for drop-in files in a set of search
// directories. `searchPath` is ordered by priority with lowest-priority
// first. That means that a drop-in file found in a latter directory will
// overwrite any drop-in file with the same name of a previous directory.
// For example, if the searchPath would equal "/var/lib/system-deploy",
// "/etc/system-deploy" then a /etc/system-deploy/<unit>/10-overwrite.conf would
// overwrite /var/lib/system-deploy/<unit>/10-overwrite.conf.
func SearchDropinFiles(unitName string, searchPath []string) ([]string, error) {
	files := make(map[string]string)

	for _, path := range searchPath {
		unitPaths := DropInSearchPaths(unitName, path)
		for _, sp := range unitPaths {
			dirFiles, err := readDir(sp)
			if os.IsNotExist(err) {
				continue
			}

			if err != nil {
				return nil, err
			}

			for _, file := range dirFiles {
				n := file.Name()
				if !file.IsDir() && strings.HasSuffix(n, DropInExt) {
					files[n] = filepath.Join(sp, n)
				}
			}
		}
	}

	// get all file names and sort them by name.
	order := make([]string, 0, len(files))
	for key := range files {
		order = append(order, key)
	}
	sort.StringSlice(order).Sort()

	// convert those file names to there actual paths
	result := make([]string, len(order))
	for idx, key := range order {
		result[idx] = files[key]
	}

	return result, nil
}

// DropInSearchPaths returns the search paths that should be checked when
// searching for task specific drop-ins. Normally, drop-ins should be placed
// in <rootDir>/<unitName>.d/<file>.conf. If the task name contains dashes
// the name is split there and used as a search path as well. In other words,
// the search path for foo-bar.task will result in the following search
// path: <rootDir>/foo-.task.d/, <rootDir>/foo-bar.task.d/. If the unitName
// contains an extension (like .task), it is used for <rootDir>/task.d/ as well.
// The returned search path is already sorted by priority where the first search
// path has lowest and the last search path has highest priority.
func DropInSearchPaths(unitName string, rootDir string) []string {
	var paths []string
	ext := filepath.Ext(unitName)
	name := strings.TrimSuffix(unitName, ext)

	// add <rootDir>/task.d
	if len(ext) > 1 { // ignore empty or dot-only extensions.
		paths = append(paths,
			filepath.Join(rootDir, strings.TrimLeft(ext, ".")+".d"),
		)
	}

	// add <rootDir>/foo-.task.d and <rootDir>/foo-bar-.task.d
	parts := strings.Split(name, "-")
	for idx := 0; idx < len(parts)-1; idx++ {
		paths = append(
			paths,
			filepath.Join(
				rootDir,
				strings.Join(parts[0:idx+1], "-")+"-"+ext+".d",
			),
		)
	}

	// add <rootDir>/foo-bar-baz.task.d
	paths = append(paths, filepath.Join(rootDir, unitName+".d"))
	return paths
}
