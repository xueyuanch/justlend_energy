package http

import (
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"justlend/internal/justlend"
	"justlend/internal/justlend/endpoints"
	"net/http"
)

func (s *Server) registerReturnResourceRouters(r *mux.Router) {
	r.Methods(http.MethodPost).Path("/return").Handler(httptransport.NewServer(
		endpoints.MakeReturnResourceEndpoint(s.service),
		decodeReturnResourceRequest,
		encodeResponse,
		s.opts...,
	))
}

// @Summary			退款.
// @Description		退款
// @Tags			交易
// @Accept			json
// @Produce			json
// @Param			receive			body		string		true	"速冲地址"
// @Param			type			body		int			true	"速冲类型0(宽带),1(能量)"
// @Param			stakePerTrx		body		int			true	"退款数量"
// @Param			privateKey		body		string		true	"扣费私钥"
// @Success			1000			{object}	justlend.ReturnResourceRL
// @Router			/return [POST]
func decodeReturnResourceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := justlend.ReturnResourceMeta{}
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return &endpoints.ReturnResourceRequest{ReturnResourceMeta: &req}, nil
}
