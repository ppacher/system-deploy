package runner

import (
	"errors"
	"sync"

	"github.com/ppacher/system-deploy/pkg/actions"
	"github.com/tevino/abool"
)

var (
	// ErrTaskExists is returned when a task or
	// task name is expect to not exist but does.
	ErrTaskExists = errors.New("task exists")

	// ErrTaskNotExists is returned when a task is
	// expected to exists but doesn't.
	ErrTaskNotExists = errors.New("task does not exist")
)

// TaskManager is responsible for managing tasks
// and implements the actions.TaskManager interface.
type TaskManager struct {
	l     sync.RWMutex
	tasks map[string]*Task
	order []string
}

// NewTaskManager returns a new task manager.
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*Task),
	}
}

// AddTask adds a new task to task manager.
func (tm *TaskManager) AddTask(name string, actions []actions.Action) error {
	t := &Task{
		actions: actions,
		name:    name,
		masked:  abool.New(),
	}

	tm.l.Lock()
	defer tm.l.Unlock()

	if _, ok := tm.tasks[name]; ok {
		return ErrTaskExists
	}

	tm.tasks[name] = t
	tm.order = append(tm.order, name)
	return nil
}

// MaskTask masks a task from execution.
func (tm *TaskManager) MaskTask(task string) error {
	t, err := tm.getTask(task)
	if err != nil {
		return err
	}

	t.mask()

	return nil
}

// UnmaskTask unmasks a task for execution.
func (tm *TaskManager) UnmaskTask(task string) error {
	t, err := tm.getTask(task)
	if err != nil {
		return err
	}

	t.unmask()

	return nil
}

// IsMasked returns true if task is masked from execution.
func (tm *TaskManager) IsMasked(task string) (bool, error) {
	t, err := tm.getTask(task)
	if err != nil {
		return false, err
	}

	return t.isMasked(), nil
}

// HasTask returns true if task exists.
func (tm *TaskManager) HasTask(task string) bool {
	_, err := tm.getTask(task)
	return err == nil
}

// IsBefore returns true if task1 is executed before task2.
func (tm *TaskManager) IsBefore(task1, task2 string) (bool, error) {
	tm.l.RLock()
	defer tm.l.RUnlock()

	t1 := -1
	t2 := -1

	for idx := range tm.order {
		if tm.order[idx] == task1 {
			t1 = idx
			continue
		}

		if tm.order[idx] == task2 {
			t2 = idx
		}
	}

	if t1 == -1 || t2 == -1 {
		return false, ErrTaskNotExists
	}

	return t1 < t2, nil
}

// IsAfter returns true if task1 is executed after task2.
func (tm *TaskManager) IsAfter(task1, task2 string) (bool, error) {
	isBefore, err := tm.IsBefore(task1, task2)
	if err != nil {
		return false, err
	}

	return !isBefore, nil
}

// getTask returns the task with the given name.
func (tm *TaskManager) getTask(name string) (*Task, error) {
	tm.l.RLock()
	defer tm.l.RUnlock()

	t, ok := tm.tasks[name]
	if !ok {
		return nil, ErrTaskNotExists
	}

	return t, nil
}

// getTaskAtIndex returns the task at the given index.
func (tm *TaskManager) getTaskAtIndex(idx int) (*Task, bool) {
	tm.l.RLock()
	defer tm.l.RUnlock()

	if idx >= len(tm.order) {
		return nil, false
	}
	name := tm.order[idx]

	task, ok := tm.tasks[name]
	return task, ok
}
