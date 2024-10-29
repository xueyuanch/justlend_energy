package derrors

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	// NotFound indicates that a requested entity was not found (HTTP 404).
	NotFound = errors.New("not found")

	// InvalidParam indicates that the input into the request is invalid in
	// some way (HTTP 400).
	InvalidParam = errors.New("invalid param")

	// Timeout means operation expired before completion.
	// For operations that change the state of the system, this error may be
	// returned even if the operation has completed successfully. For
	// example, a successful response from a server could have been delayed
	// long enough for the deadline to expire.
	Timeout = errors.New("timeout")

	// Unknown error.
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	Unknown = errors.New("unknown")

	// Internal server error.
	// error raised while server encountered an internal error or
	// misconfiguration and was unable to complete user request.
	Internal = errors.New("internal")

	// Forbidden indicates the caller does not have permission to
	// execute the specified operation. It must not be used for rejections
	// caused by exhausting some resource (use ResourceExhausted
	// instead for those errors). It must not be used if the caller
	// cannot be identified (use Unauthenticated instead for those errors).
	Forbidden = errors.New("forbidden")

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	Unauthenticated = errors.New("unauthenticated")

	// InconsistentData indicate the data is inconsistent inertly or inconsistent
	// with system data. This can be caused by the client-side cache being out
	// of sync with the server data.
	InconsistentData = errors.New("inconsistent data")

	// Unavailable indicates the service is currently unavailable.
	// This is a most likely a transient condition and may be corrected
	// by retrying with a backoff. Note that it is not always safe to retry
	// non-idempotent operations.
	Unavailable = errors.New("unavailable")

	// Duplicate represent a duplicate operation/request, i.e. non-unique
	// login or other request operations/requests.
	Duplicate = errors.New("duplicate")
)

// Add adds context to the error.
// The result cannot be unwrapped to recover the original error.
// It does nothing when *err == nil.
// Example:
// defer derrors.Add(&err, "copy(%s, %s)", src, dst)
// See Wrap for an equivalent function that allows
// the result to be unwrapped.
func Add(err *error, format string, args ...interface{}) {
	if *err == nil {
		// Abort if the original error is nil.
		return
	}
	*err = fmt.Errorf("%s: %v", fmt.Sprintf(format, args...), *err)
}

// Wrap adds context to the error and allows
// unwrapping the result to recover the original error.
//
// Example:
//
//	defer derrors.Wrap(&err, "copy(%s, %s)", src, dst)
//
// See Add for an equivalent function that does not allow
// the result to be unwrapped.
func Wrap(err *error, format string, args ...interface{}) {
	if *err == nil {
		// Abort if the original error is nil.
		return
	}
	// The verb %w is used here that we can recover the
	// original error in the future.
	*err = fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), *err)
}

// WrapStack is like Wrap, but adds a stack trace if there isn't one already.
func WrapStack(err *error, format string, args ...interface{}) {
	if *err == nil {
		// Abort if the original error is nil.
		return
	}
	// Wrap error with StackError which will capturing the current stack
	// trace if the original error is not of type *StackError
	if se := (*stackError)(nil); !errors.As(*err, &se) {
		*err = NewStackError(*err)
	}
	Wrap(err, format, args...)
}

// Unwrap returns the result of calling the Unwrap method on err, if errs
// type contains an Unwrap method returning error. Otherwise, Unwrap
// returns original error.
func Unwrap(err error) error {
	for {
		if u, ok := err.(interface{ Unwrap() error }); ok {
			err = u.Unwrap()
		} else {
			return err
		}
	}
}

// IsNotFound is helps to check if a particular error is NotFound failure.
func IsNotFound(err error) bool { return errors.Is(Unwrap(err), NotFound) }
