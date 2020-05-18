package actions

import (
	"context"
	"strings"
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

// OptionSpecs returns a map using lower-case option names
// as the key.
func (plg *Plugin) OptionSpecs() map[string]deploy.OptionSpec {
	m := make(map[string]deploy.OptionSpec, len(plg.Options))

	for _, opt := range plg.Options {
		m[strings.ToLower(opt.Name)] = opt
	}

	return m
}

// Action describes a generic action that is capable of performing
// a single task. Actions are grouped into tasks and are executed in
// the order they are defined inside a task. Action implementations
// are encouraged to embed the Base struct defined below.
type Action interface {
	// Name should return a name for the action.
	Name() string

	// SetLogger is called before Setup and configures the logger
	SetLogger(l Logger)

	// SetTask configures the deploy task.
	SetTask(t deploy.Task)
}

// Preparer describes the interface that actions can implement
// if they want to prepare execution or perform other actions
// during the preperation phase.
type Preparer interface {
	// Prepare should prepare the action and return
	// whether or not the task should be executed or not.
	Prepare(ExecGraph) error
}

// Executor describes the interface that actions can implement if
// they want to run during the execution phase.
type Executor interface {
	// Execute actually performs the action. The returned
	// boolean should be set to true if the action
	// actually did some modifications.
	Execute(ctx context.Context) (bool, error)
}

var (
	actionsLock sync.RWMutex
	actions     map[string]*Plugin
)
