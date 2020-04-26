package deploy

import (
	"errors"
	"fmt"
	"io"
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
	Description string

	// StartMasked is set to true if this task is disabled (masked)
	// by default.
	StartMasked bool

	// Sections holds the tasks sections.
	Sections []unit.Section
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
				return nil, fmt.Errorf("failed to decode [Task] section")
			}

			task.Sections = append(sections[:idx], sections[idx+1:]...)

			break
		}
	}

	if len(task.Sections) == 0 {
		return nil, fmt.Errorf("deploy task does not define any actions")
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
	}

	startMasked, err := section.GetBool("StartMasked")
	if err != nil && err != unit.ErrOptionNotSet {
		return fmt.Errorf("error in option 'StartMasked': %w", err)
	}

	task.Description = description
	task.StartMasked = startMasked

	return nil
}
