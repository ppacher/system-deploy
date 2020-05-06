package runner

import (
	"context"

	"github.com/ppacher/system-deploy/pkg/actions"
	"github.com/tevino/abool"
)

type Task struct {
	actions []actions.Action

	name     string
	masked   *abool.AtomicBool
	disabled *abool.AtomicBool
}

// mask the task from execution. If t is a nil task mask is a no-op.
func (t *Task) mask() {
	if t == nil {
		return
	}
	t.masked.Set()
}

// unmask the task for execution. If t is a nil task unmask is a no-op.
func (t *Task) unmask() {
	if t == nil {
		return
	}
	t.masked.UnSet()
}

// isMasked returns true if t is masked from execution or t is a
// nil task.
func (t *Task) isMasked() bool {
	if t == nil {
		return true
	}

	return t.masked.IsSet()
}

// Prepare calls the perpare method of each action defined
// in the task.
func (t *Task) Prepare(graph actions.ExecGraph) error {
	for _, a := range t.actions {
		if p, ok := a.(actions.Preparer); ok {
			if err := p.Prepare(graph); err != nil {
				return err
			}
		}
	}

	return nil
}

// Execute executes all actions of the task in the order they are defined.
// It returns true if any of the actions returned true and aborts on the
// first error encountered.
func (t *Task) Execute(ctx context.Context, log actions.Logger) (bool, error) {
	var changed bool
	for _, a := range t.actions {
		if r, ok := a.(actions.Executor); ok {
			c, err := r.Execute(ctx)
			if err != nil {
				return false, err
			}

			if !changed {
				changed = c
			}
		}
	}

	return changed, nil
}
