package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"justlend/internal/justlend"
)

type ReturnResourceRequest struct {
	*justlend.ReturnResourceMeta
}

func MakeReturnResourceEndpoint(s justlend.Service) endpoint.Endpoint {
	return Sentry(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*ReturnResourceRequest)
		return NewResponse(s.ReturnResource(ctx, req.ReturnResourceMeta)), nil
	})
}
