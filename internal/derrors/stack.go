package derrors

import "runtime"

// stackError wraps an error and adds a stack trace.
type stackError struct {
	Stack []byte
	err   error
}

// NewStackError returns a StackError, capturing a stack trace.
func NewStackError(err error) *stackError {
	// Limit the stack trace to 16K. Same value used in the error reporting client,
	// https://pkg.go.dev/cloud.google.com/go/errorreporting.
	var buf [16 * 1024]byte
	n := runtime.Stack(buf[:], false)
	return &stackError{
		err:   err,
		Stack: buf[:n],
	}
}

func (e *stackError) Error() string {
	return e.err.Error() // ignore the stack
}

func (e *stackError) Unwrap() error {
	return e.err
}
