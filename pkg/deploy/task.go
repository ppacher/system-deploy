package deploy

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ppacher/system-deploy/pkg/unit"
)

// Task defines a deploy task.
type Task struct {
	// FileName is the name of the file that describes this task.
	FileName string

	// Directory holds the directory of the task.
	Directory string

	// Description is the tasks description.
	Description *string

	// StartMasked is set to true if this task is disabled (masked)
	// by default.
	StartMasked *bool

	// Disabled can be set to true to disable a task permanently.
	Disabled *bool

	// Sections holds the tasks sections.
	Sections []unit.Section
}

// DecodeFile is like Decode but reads the task from
// filePath.
func DecodeFile(filePath string) (*Task, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Decode(filePath, f)
}

// Decode decodes a deploy task from r and uses the basename
// of name as the task's name.
func Decode(filePath string, r io.Reader) (*Task, error) {
	sections, err := unit.Deserialize(r)
	if err != nil {
		return nil, err
	}

	task := &Task{
		FileName:  filepath.Base(filePath),
		Directory: filepath.Dir(filePath),
	}

	for idx, sec := range sections {
		if strings.ToLower(sec.Name) == "task" {
			if err := decodeMetaData(sec, task); err != nil {
				return nil, ErrInvalidTaskSection
			}

			task.Sections = append(sections[:idx], sections[idx+1:]...)

			break
		}
	}

	if len(task.Sections) == 0 {
		return nil, ErrNoSections
	}

	return task, nil
}

func decodeMetaData(section unit.Section, task *Task) error {
	if strings.ToLower(section.Name) != "task" {
		return errors.New("invalid section name")
	}

	description, err := section.GetString("Description")
	if err != nil && err != unit.ErrOptionNotSet {
		return fmt.Errorf("error in option 'Description': %w", err)
	} else if err == nil {
		task.Description = &description
	}

	startMasked, err := section.GetBool("StartMasked")
	if err != nil && err != unit.ErrOptionNotSet {
		return fmt.Errorf("error in option 'StartMasked': %w", err)
	} else if err == nil {
		task.StartMasked = &startMasked
	}

	disabled, err := section.GetBool("Disabled")
	if err != nil && err != unit.ErrOptionNotSet {
		return fmt.Errorf("error in option 'Disabled': %w", err)
	} else if err == nil {
		task.Disabled = &disabled
	}

	return nil
}

// Clone creates a deep copy of t.
func (t *Task) Clone() *Task {
	n := &Task{
		FileName:    t.FileName,
		Directory:   t.Directory,
		Description: t.Description,
		StartMasked: t.StartMasked,
		Disabled:    t.Disabled,
	}

	if len(t.Sections) > 0 {
		n.Sections = make([]unit.Section, len(t.Sections))
		for idx, s := range t.Sections {
			n.Sections[idx] = unit.Section{
				Name:    s.Name,
				Options: make(unit.Options, len(s.Options)),
			}

			for optIdx, opt := range s.Options {
				n.Sections[idx].Options[optIdx] = unit.Option{
					Name:  opt.Name,
					Value: opt.Value,
				}
			}
		}
	} else if t.Sections != nil {
		// make sure we also have an empty slice
		n.Sections = make([]unit.Section, 0)
	}

	return n
}
