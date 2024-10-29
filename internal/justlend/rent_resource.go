package justlend

import (
	"context"
	"justlend/internal"
	"justlend/internal/derrors"
	"justlend/internal/protos/core"
)

type RentResourceMeta struct {
	Receive    string            `json:"receive"`
	Type       core.ResourceCode `json:"type"`
	Amount     int64             `json:"amount"`
	PrivateKey string            `json:"privateKey"`
}

func (m *RentResourceMeta) Conform(_ context.Context) error {
	switch {
	case !internal.IsValidAddress(m.Receive):
		return derrors.InvalidParam
	case !internal.Contains(m.Type, core.ResourceCode_BANDWIDTH, core.ResourceCode_ENERGY):
		return derrors.InvalidParam
	case m.Amount <= 0:
		return derrors.InvalidParam
	case len(m.PrivateKey) != 64:
		return derrors.InvalidParam
	default:
		return nil
	}
}

type RentResourceRL struct {
	TxId        string `json:"txId"`
	StakePerTrx int64  `json:"stakePerTrx"`
}

var (
	_ internal.Conformer = (*RentResourceMeta)(nil)
)

type RentResourceService interface {
	RentResource(ctx context.Context, req *RentResourceMeta) (*RentResourceRL, error)
}
