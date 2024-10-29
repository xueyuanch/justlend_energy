package derrors

import (
	"github.com/pkg/errors"
	"net/http"
)

const (
	// The following code define a unique identifier for the error, The code should
	// not match the response status code. Instead, it should be an error code unique
	// to our application.
	// Generally, there is no convention for the codes, expect that it be unique
	// numerical code across the application services.

	// cSuccess means operation success.
	cSuccess = 1000

	// 40xx codes represents the generic errors across all
	// application services.
	cInternal         = 4000
	cNotFound         = 4001
	cInvalidArgument  = 4002
	cTimeout          = 4003
	cUnknown          = 4004
	cForbidden        = 4005
	cUnauthenticated  = 4006
	cInconsistentData = 4007
	cUnavailable      = 4008
	cDuplicate        = 4009
)

var codes = []struct {
	err     error
	code    int    // Machine-readable message.
	HTTP    int    // HTTP status code.
	message string // A brief human-readable message.
}{

	{
		nil,
		cSuccess,
		http.StatusOK, // Standard http status codes.
		"",
	}, // Zero value for code and message is used for nil-error.
	{
		Internal,
		cInternal,
		http.StatusInternalServerError,
		"服务器开小差了, 请稍后再试",
	}, // Default code for undefined errors.
	{
		NotFound,
		cNotFound,
		http.StatusNotFound,
		"资源未找到",
	},
	{
		InvalidParam,
		cInvalidArgument,
		http.StatusBadRequest,
		"无效参数",
	},
	{
		InconsistentData,
		cInconsistentData,
		http.StatusBadRequest,
		"数据已过期, 请刷新再试",
	},
	{
		Forbidden,
		cForbidden,
		http.StatusForbidden,
		"禁止访问",
	},
	{
		Unauthenticated,
		cUnauthenticated,
		http.StatusUnauthorized,
		"操作未授权",
	},
	{
		Unavailable,
		cUnavailable,
		http.StatusServiceUnavailable,
		"资源不可用",
	},
	{
		Timeout,
		cTimeout,
		http.StatusRequestTimeout,
		"请求超时",
	},
	{
		Duplicate,
		cDuplicate,
		http.StatusConflict,
		"重复操作",
	},
	{
		Unknown,
		cUnknown,
		http.StatusBadRequest,
		"未知操作",
	},
}

// ToCode returns a unique code and message corresponding to err.
// If err is nil, it returns the zero values.
func ToCode(err error) (int, string) {
	originalErr := Unwrap(err)
	for _, c := range codes {
		if errors.Is(originalErr, c.err) {
			return c.code, c.message
		}
	}
	if fer, ok := originalErr.(failure); ok {
		// fer represents an application-specific failure error. Application errors
		// can be unwrapped by the caller to extract out the string message.
		//
		// Any non-application error (such as a disk error) should be reported as an
		// Internal error and the human user should only see "Internal error" as the
		// message. These low-level internal error details should only be logged and
		// reported to the operator of the application (not the end user).
		return cInvalidArgument, fer.failure()
	}
	// cInternal returned if error not defined in
	// codes which means an unexpected error occurs, and we should
	// inspect what went wrong.
	return cInternal, "服务器开小差了, 请稍后再试"
}

func ToError(code int) error {
	for _, c := range codes {
		if c.code == code {
			return c.err
		}
	}
	return nil
}
