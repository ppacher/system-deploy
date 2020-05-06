package actions

import "github.com/ppacher/system-deploy/pkg/deploy"

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
