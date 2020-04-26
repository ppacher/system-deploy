package onchange

import (
	"context"
	"fmt"

	"github.com/ppacher/system-deploy/pkg/actions"
	"github.com/ppacher/system-deploy/pkg/deploy"
	"github.com/ppacher/system-deploy/pkg/unit"
	"github.com/ppacher/system-deploy/pkg/utils"
)

func init() {
	actions.MustRegister(actions.Plugin{
		Name:        "OnChange",
		Description: "Perform post operations after the current task.",
		Setup:       setupAction,
		Author:      "Patrick Pacher <patrick.pacher@gmail.com>",
		Website:     "https://github.com/ppacher/system-deploy",
		Options: []deploy.OptionSpec{
			{
				Name:        "Run",
				Type:        deploy.StringSliceType,
				Description: "Run a command. May be specified multiple times. Note that errors are only logged and don't abort subsequent tasks. Use Unmask for more control",
			},
			{
				Name:        "Unmask",
				Type:        deploy.StringSliceType,
				Description: "Unmask a task. May be specified multiple times.",
			},
		},
	})
}

func setupAction(task deploy.Task, sec unit.Section) (actions.Action, error) {
	return &action{
		task: task,
		sec:  sec,
	}, nil
}

type action struct {
	actions.Base

	task deploy.Task
	sec  unit.Section
}

func (a *action) Name() string {
	return "OnChange"
}

func (a *action) Prepare(graph actions.ExecGraph) error {
	err := a.forEachStringValue("Run", func(value string) error {
		return a.runOnChange(graph, func(ctx context.Context) {
			utils.ExecCommand(ctx, a.task.Directory, value, nil)
		})
	})
	if err != nil {
		return err
	}

	err = a.forEachStringValue("Unmask", func(value string) error {
		if !graph.HasTask(value) {
			return fmt.Errorf("unknown task %s", value)
		}
		return a.runOnChange(graph, func(ctx context.Context) {
			a.Debugf("Unmasking task %s", value)
			if err := graph.UnmaskTask(value); err != nil {
				a.Warnf("failed to unmask task %q: %s", value, err)
			}
		})
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *action) runOnChange(graph actions.ExecGraph, fn func(context.Context)) error {
	return graph.RunAfter(a.task.FileName, func(ctx context.Context, _ string, update bool, err error) {
		if err != nil {
			return
		}

		if !update {
			return
		}

		fn(ctx)
	})
}

func (a *action) forEachStringValue(key string, fn func(string) error) error {
	values := a.sec.GetStringSlice(key)

	for _, v := range values {
		if err := fn(v); err != nil {
			return err
		}
	}

	return nil
}
