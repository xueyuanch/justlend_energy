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

func (s *Server) registerRentResourceRouters(r *mux.Router) {
	r.Methods(http.MethodPost).Path("/rent").Handler(httptransport.NewServer(
		endpoints.MakeRentResourceEndpoint(s.service),
		decodeRentResourceRequest,
		encodeResponse,
		s.opts...,
	))
}

// @Summary			租用.
// @Description		租用
// @Tags			交易
// @Accept			json
// @Produce			json
// @Param			receive			body		string		true	"速冲地址"
// @Param			type			body		int			true	"速冲类型0(宽带),1(能量)"
// @Param			amount			body		int			true	"速冲数量"
// @Param			privateKey		body		string		true	"扣费私钥"
// @Success			1000			{object}	justlend.RentResourceRL
// @Router			/rent [POST]
func decodeRentResourceRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := justlend.RentResourceMeta{}
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return &endpoints.RentResourceRequest{RentResourceMeta: &req}, nil
}
