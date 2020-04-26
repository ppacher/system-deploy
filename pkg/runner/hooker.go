package runner

import (
	"context"
	"sync"

	"github.com/ppacher/system-deploy/pkg/actions"
)

// Hooker allows to register before and after task execution
// hooks. It implements the actions.Hooker interface.
type Hooker struct {
	l         sync.RWMutex
	preHooks  map[string][]actions.BeforeTaskFunc
	postHooks map[string][]actions.AfterTaskFunc
}

// NewHooker creates and returns a new hooker.
func NewHooker() *Hooker {
	return &Hooker{
		preHooks:  make(map[string][]actions.BeforeTaskFunc),
		postHooks: make(map[string][]actions.AfterTaskFunc),
	}
}

// RunBefore registeres fn to be executed before task.
func (h *Hooker) RunBefore(task string, fn actions.BeforeTaskFunc) error {
	h.l.Lock()
	defer h.l.Unlock()

	hooks := h.preHooks[task]
	hooks = append(hooks, fn)

	h.preHooks[task] = hooks
	return nil
}

// RunAfter registeres fn to be executed after task.
func (h *Hooker) RunAfter(task string, fn actions.AfterTaskFunc) error {
	h.l.Lock()
	defer h.l.Unlock()

	hooks := h.postHooks[task]
	hooks = append(hooks, fn)

	h.postHooks[task] = hooks
	return nil
}

// ExecuteBefore executes all BeforeFunc that have been registered for task.
// It aborts as soon as a function returns an error.
func (h *Hooker) ExecuteBefore(ctx context.Context, task string) (context.Context, error) {
	h.l.RLock()
	defer h.l.RUnlock()

	var err error
	for _, fn := range h.preHooks[task] {
		ctx, err = fn(ctx, task)
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

// ExecuteAfter executes all AfterFunc that have been registered for task.
func (h *Hooker) ExecuteAfter(ctx context.Context, task string, actionPerformed bool, errResult error) {
	h.l.RLock()
	defer h.l.RUnlock()

	for _, fn := range h.postHooks[task] {
		fn(ctx, task, actionPerformed, errResult)
	}
}
