package justlend

import (
	"context"
	"justlend/internal"
	"justlend/internal/derrors"
	"justlend/internal/protos/core"
)

type ReturnResourceMeta struct {
	Receive     string            `json:"receive"`
	Type        core.ResourceCode `json:"type"`
	StakePerTrx int64             `json:"stakePerTrx"`
	PrivateKey  string            `json:"privateKey"`
}

func (m *ReturnResourceMeta) Conform(_ context.Context) error {
	switch {
	case !internal.IsValidAddress(m.Receive):
		return derrors.InvalidParam
	case !internal.Contains(m.Type, core.ResourceCode_BANDWIDTH, core.ResourceCode_ENERGY):
		return derrors.InvalidParam
	case m.StakePerTrx <= 0:
		return derrors.InvalidParam
	case len(m.PrivateKey) != 64:
		return derrors.InvalidParam
	default:
		return nil
	}
}

type ReturnResourceRL struct {
	TxId string `json:"txId"`
}

var (
	_ internal.Conformer = (*ReturnResourceMeta)(nil)
)

type ReturnResourceService interface {
	ReturnResource(ctx context.Context, req *ReturnResourceMeta) (*ReturnResourceRL, error)
}
