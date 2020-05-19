package runner

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/ppacher/system-deploy/pkg/actions"
	"github.com/ppacher/system-deploy/pkg/deploy"
)

// Runner executes a set of targets in order and aborts
// on the first error.
type Runner struct {
	*TaskManager
	*Hooker

	l actions.Logger
}

// NewRunner creates a new runner for the given targets.
func NewRunner(l actions.Logger, targets []deploy.Task) (*Runner, error) {
	r := &Runner{
		TaskManager: NewTaskManager(l),
		Hooker:      NewHooker(),
		l:           l,
	}

	for _, target := range targets {
		if err := r.AddTask(target.FileName, target); err != nil {
			return nil, fmt.Errorf("failed to add target %s: %w", target.FileName, err)
		}
	}

	return r, nil
}

// Deploy runs all deploy targets and aborts and returns
// the first error encountered.
func (r *Runner) Deploy(ctx context.Context) error {
	iter := &taskIter{
		tm: r.TaskManager,
	}

	r.inPrepare.Set()
	for iter.Next() {
		r.l.Debugf("Preparing task %q", iter.Name())
		if err := iter.Task().Prepare(r); err != nil {
			return fmt.Errorf("failed to perpare target %s: %w", iter.Name(), err)
		}
	}
	r.inPrepare.UnSet()

	iter.Reset()

	bold := color.New(color.Bold)

	r.inExec.Set()
	defer r.inExec.UnSet()

	for iter.Next() {
		if iter.IsMasked() {
			r.l.Debugf("skipping masked target %s", iter.Name())
			continue
		}

		task := iter.Task()
		name := iter.Name()

		taskContext, err := r.ExecuteBefore(ctx, name)
		if err != nil {
			return err
		}

		r.l.Debugf("Starting task %s", bold.Sprint(name))
		res, err := task.Execute(taskContext, r.l)

		r.ExecuteAfter(taskContext, name, res, err)

		if err != nil {
			r.l.Warnf("%s: %s", color.New(color.BgRed, color.FgWhite).Sprint("FAIL"), err.Error())
			return err
		}
		resStr := "pristine"

		if res {
			resStr = color.New(color.FgHiGreen, color.Bold).Sprint("updated")
		}
		r.l.Infof("%s: %s", bold.Sprintf("%-30v", name), resStr)
	}

	return nil
}
