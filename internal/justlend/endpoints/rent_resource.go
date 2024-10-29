package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"justlend/internal/justlend"
)

type RentResourceRequest struct {
	*justlend.RentResourceMeta
}

func MakeRentResourceEndpoint(s justlend.Service) endpoint.Endpoint {
	return Sentry(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*RentResourceRequest)
		return NewResponse(s.RentResource(ctx, req.RentResourceMeta)), nil
	})
}
