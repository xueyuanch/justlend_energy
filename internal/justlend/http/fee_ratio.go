package http

import (
	"context"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"justlend/internal/justlend"
	"justlend/internal/justlend/endpoints"
	"justlend/internal/protos/core"
	"net/http"
)

func (s *Server) registerFeeRatioRouters(r *mux.Router) {
	r.Methods(http.MethodGet).Path("/fee").Handler(httptransport.NewServer(
		endpoints.MakeFeeRatioEndpoint(s.service),
		decodeFeeRatioRequest,
		encodeResponse,
		s.opts...,
	))
}

// @Summary			计算费用.
// @Description		根据参数计算费用
// @Tags			费用计算
// @Accept			json
// @Produce			json
// @Param			energy			query		int		true	"需要速冲的数量"
// @Param			privateKey		query		string	true	"私钥"
// @Param			type			query		int32	true	"0(宽带),1(能量)"
// @Success			1000			{object}	justlend.FeeRatioRL
// @Router			/fee [GET]
func decodeFeeRatioRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return &endpoints.FeeRatioRequest{
		FeeRatioMeta: &justlend.FeeRatioMeta{
			Energy:     safeExtractQueryInt(r, "energy"),
			PrivateKey: safeExtractQueryString(r, "privateKey"),
			Type:       core.ResourceCode(safeExtractQueryInt(r, "type")),
		},
	}, nil
}
