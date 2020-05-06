package actions

import "context"

// PreRunFunc is executed by the execution graph before
// the actual tasks are executed. The returned context will
// be passed to each task's run method.
type PreRunFunc func(context.Context) (context.Context, error)

// PostRunFunc is executed by the execution graph after
// the actual tasks are executed. It receives an overall
// success status indicating if _all_ tasks executed
// successfully.
type PostRunFunc func(ctx context.Context, success bool)

// TaskManager allows to mask certain tasks from execution.
type TaskManager interface {
	// MaskTask masks a task from execution.
	MaskTask(task string) error

	// UnmaskTask unmasks a task for execution.
	UnmaskTask(task string) error

	// DisableTask disables a task. Once disabled, a task cannot be
	// enabled anymore. This should only be used if it's likely that
	// a given task cannot succeed under certain conditions like
	// an unsupported operating system or package manager. Once a
	// task is disabled it won't run in the perperation and execution
	// phase. DisableTask can only be called during Prepare() and will
	// return an error during Run(). Note that depending on the order,
	// a task's Perpare() function might have already been executed when
	// the task gets disabled.
	DisableTask(task string) error

	// IsMasked returns true if a task is masked from execution.
	IsMasked(task string) (bool, error)

	// HasTask checks if task exists.
	HasTask(task string) bool

	// IsBefore returns true if task1 would executed before task2
	IsBefore(task1, task2 string) (bool, error)

	// IsAfter returns true if task1 would be executed after task2
	IsAfter(task1, task2 string) (bool, error)
}

// BeforeTaskFunc is executed before a certain task is executed.
type BeforeTaskFunc func(ctx context.Context, task string) (context.Context, error)

// AfterTaskFunc is executed after a certain task has been
// executed. It receives the error (or nil) returned by the
// task.
type AfterTaskFunc func(ctx context.Context, task string, actionPerformed bool, err error)

// Hooker allows to place hooks that are run before or
// after a certain task.
type Hooker interface {
	// RunBefore executes fn before task. Note that fn is only
	// called if task is really going to be executed and not
	// masked in the TaskManager.
	RunBefore(task string, fn BeforeTaskFunc) error

	// RunAfter executes fn after task has been executed.
	// Note that fn is only called if task was really
	// executed and not masked in the TaskManager
	RunAfter(task string, fn AfterTaskFunc) error
}

// ExecGraph defines the execution graph for deploy tasks.
type ExecGraph interface {
	TaskManager
	Hooker

	// AddPreRun registers a pre-run function that is executed
	// before the execution graph starts.
	//AddPreRun(fn PreRunFunc)

	// AddPostRun registeres a post-run function that is executed
	// afterh the execution graph finished.
	//AddPostRun(fn PostRunFunc)
}
