package actions

import (
	"github.com/sirupsen/logrus"
)

// Logger defines the interface that is passed to actions
// and that should be used to keep the user updated on
// what's currently happening.
type Logger interface {
	// Progress can be used to report the current progress
	// if available.
	Progress(value float64, msg string)

	// Infof logs an informative message.
	Infof(fmt string, args ...interface{})

	// Debugf logs a debug message mainly meant for
	// developers.
	Debugf(fmt string, args ...interface{})

	// Warnf logs a warning message.
	Warnf(fmt string, args ...interface{})
}

type simpleLogger struct {
	*logrus.Logger
}

func (sl *simpleLogger) Progress(value float64, msg string) {
	sl.Infof("Progress %.2f%%: %s", value, msg)
}

// NewLogger returns a new logger
func NewLogger() Logger {
	l := &simpleLogger{
		Logger: logrus.New(),
	}

	l.SetLevel(logrus.GetLevel())

	return l
}
