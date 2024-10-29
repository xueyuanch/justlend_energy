package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"justlend/internal/justlend"
)

type FeeRatioRequest struct {
	*justlend.FeeRatioMeta
}

func MakeFeeRatioEndpoint(s justlend.Service) endpoint.Endpoint {
	return Sentry(func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*FeeRatioRequest)
		return NewResponse(s.FeeRatio(ctx, req.FeeRatioMeta)), nil
	})
}
