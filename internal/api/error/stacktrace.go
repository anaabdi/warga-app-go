package error

import (
	"fmt"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// StackTrace wrap err with stack trace
func StackTrace(err error) error {
	return stackTraceMessage(err, nil)
}

// StackTraceMessage wrap err with stack trace and additional err message
func StackTraceMessage(err error, message *string) error {
	return stackTraceMessage(err, message)
}

func stackTraceMessage(err error, message *string) error {
	if err == nil {
		return nil
	}
	st := errors.WithStack(err).(stackTracer).StackTrace()[2:]

	if message == nil {
		return fmt.Errorf("%w%+v", err, st)
	}
	return fmt.Errorf("%v: %w%+v", *message, err, st)
}
