package actions

import (
	"context"
	"sync"

	"github.com/ppacher/system-deploy/pkg/deploy"
	"github.com/ppacher/system-deploy/pkg/unit"
)

// ActionFunc performs a custom action and returns either success or failure.
type ActionFunc func(ctx context.Context) error

// PostActionFunc is executed after each primary action has been performed.
// PostActionFuncs are only executed after all deploy tasks have been
// executed.
type PostActionFunc func(ctx context.Context) error

// SetupFunc should return a new action instance.
type SetupFunc func(deploy.Task, unit.Section) (Action, error)

// Plugin describes a deploy plugin.
type Plugin struct {
	// Name is the name of the plugin and used
	// to find matching sections.
	Name string

	// Description is a human readable description of
	// the plugins purpose. Description should be a
	// single short line. For more help text about
	// the plugins purpose and functioning use
	// the Help section list.
	Description string

	// Setup creates a new action base on deploy options.
	Setup SetupFunc

	// Help may contain additional help sections.
	Help []deploy.Section

	// Example may contain an example task.
	Example string

	// Options defines all supported deploy options.
	Options []deploy.OptionSpec

	// Author may hold the name of the plugin author.
	Author string

	// Website may hold the name of the plugin website.
	Website string
}

type Action interface {
	// Name should return a name for the action.
	Name() string

	// Prepare should prepare the action and return
	// whether or not the task should be executed or not.
	Prepare(ExecGraph) error

	// SetLogger is called before Setup and configures the logger
	SetLogger(l Logger)

	// SetTask configures the deploy task.
	SetTask(t deploy.Task)
}

type Runner interface {
	// Run actually performs the action. The returned
	// boolean should be set to true if the action
	// actually did some modifications.
	Run(ctx context.Context) (bool, error)
}

// Base provides a base action and is meant to be embedded into
// real action implementations.
type Base struct {
	Logger
	deploy.Task
}

// SetLogger configures the logger to use and implements
// SetLogger from actions.Action.
func (b *Base) SetLogger(l Logger) {
	b.Logger = l
}

// SetTask configures the deploy task and implements
// SetTask from actions.Action.
func (b *Base) SetTask(t deploy.Task) {
	b.Task = t
}

var (
	actionsLock sync.RWMutex
	actions     map[string]*Plugin
)
