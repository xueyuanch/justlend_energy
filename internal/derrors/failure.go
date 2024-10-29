package derrors

import "fmt"

// failure describe a domain logic failure error which expected to
// implement both error and failure.
type failure interface{ failure() string }

// failureError represents an application-specific error and it will be
// described as a `InvalidParam` error and response with a HTTP 400 status
// to consumer. Application errors can be unwrapped by the caller
// to extract out the message.
//
// Any non-application error (such as a disk error) should be reported as an
// `Internal` error and the human user should only see "Internal error" as the
// message. These low-level internal error details should only be logged and
// reported to the operator of the application (not the end user).
//
// The best practice is use `%w` verb wrapping the origin error add
// human-readable context to the underlying error and record the file
// and line that the error occurred.
type failureError string

// Implementation of errors.Error.
func (e failureError) Error() string {
	return e.String()
}

// Implementation of fmt.Stringer.
func (e failureError) String() string {
	return string(e)
}

// Failure returns the failure message which describe a
// Human-readable error message.

func (e failureError) failure() string {
	return e.String()
}

var (
	_ error        = (*failureError)(nil)
	_ fmt.Stringer = (*failureError)(nil)
)

// NewFailure annotates failureError's message with the format specifier.
func NewFailure(format string, args ...interface{}) error {
	return failureError(fmt.Sprintf(format, args...))
}
