package runner

import (
	"sync"
)

type taskIter struct {
	sync.Mutex
	idx  int
	task *Task
	tm   *TaskManager
}

func (iter *taskIter) Reset() {
	iter.Lock()
	defer iter.Unlock()

	iter.idx = 0
	iter.task = nil
}

// Next moves the taskIterator to the next task. It returns
// true if more tasks are available or false if it was
// the last task. Next is meant to be called in a for-loop.
func (iter *taskIter) Next() bool {
	iter.Lock()
	defer iter.Unlock()

	hasNext := false
	iter.task, hasNext = iter.tm.getTaskAtIndex(iter.idx)
	if hasNext {
		iter.idx++
		return true
	}

	return false
}

// IsMasked returns true if the current task is masked.
func (iter *taskIter) IsMasked() bool {
	iter.Lock()
	defer iter.Unlock()

	if iter.task == nil {
		return false
	}

	return iter.task.isMasked()
}

// Name returns the action of the current task.
func (iter *taskIter) Name() string {
	iter.Lock()
	defer iter.Unlock()

	if iter.task == nil {
		return ""
	}

	return iter.task.name
}

// Action returns the action of the current task.
func (iter *taskIter) Task() *Task {
	iter.Lock()
	defer iter.Unlock()

	return iter.task
}
