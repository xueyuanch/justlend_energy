package endpoints

import (
	"context"
	"justlend/internal"

	"github.com/go-kit/kit/endpoint"
)

// Sentry is an endpoints sentry middleware between the request and actual service
// handler, if the request parameter implements the internal.Conformer interface, call
// the conform method of the interface, which should include the basic validation of
// the request data, returns any error indicates that the data does not conform our
// expectations and will abort the subsequent chain handlers.
func Sentry(e endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if cfr, ok := req.(internal.Conformer); ok {
			// We simply return the domain failure error.
			if err := cfr.Conform(ctx); err != nil {
				return NewErrResponse(err), nil
			}
		}
		return e(ctx, req)
	}
}
