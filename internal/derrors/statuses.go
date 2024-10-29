package derrors

import (
	"github.com/pkg/errors"
	gcodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// ToHTTPSta returns an HTTP status code corresponding to err.
// If err is nil, it returns StatusOK.
func ToHTTPSta(err error) int {

	originalErr := Unwrap(err)
	for _, e := range codes {
		if errors.Is(originalErr, e.err) {
			return e.HTTP
		}
	}
	if _, ok := originalErr.(failure); ok {
		// failure implementer represents an application-specific failure
		// error, and it will be treated as a bad request.
		return http.StatusBadRequest
	}
	// StatusInternalServerError returned if error not defined in
	// codes which means an unexpected error occurs, and we should
	// inspect what went wrong.
	return http.StatusInternalServerError
}

// TogRPCSta returns a gRPC status code corresponding to err.
// If err is nil, it returns `codes.OK`.
func TogRPCSta(err error) *status.Status {
	originalErr := Unwrap(err)
	for _, e := range codes {
		if errors.Is(originalErr, e.err) {
			return status.New(gcodes.Code(e.code), e.message)
		}
	}
	if fail, ok := originalErr.(failure); ok {
		// failure implementer represents an application-specific failure
		// error, and it will be treated as a bad request.
		return status.New(gcodes.Code(cInvalidArgument), fail.failure())
	}
	// StatusInternalServerError returned if error not defined in
	// codes which means an unexpected error occurs, and we should
	// inspect what went wrong.
	return status.New(gcodes.Code(cInternal), "服务器异常, 请稍后再试")
}
