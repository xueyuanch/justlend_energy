package justlend

import (
	"context"
	"github.com/shopspring/decimal"
	"justlend/internal"
	"justlend/internal/derrors"
	"justlend/internal/protos/core"
)

type FeeRatioMeta struct {
	Energy     int64
	PrivateKey string
	Type       core.ResourceCode
}

func (m *FeeRatioMeta) Conform(_ context.Context) error {
	switch {
	case m.Energy <= 0:
		return derrors.InvalidParam
	case len(m.PrivateKey) != 64:
		return derrors.InvalidParam
	case !internal.Contains(m.Type,
		core.ResourceCode_BANDWIDTH,
		core.ResourceCode_ENERGY,
	):
		return derrors.InvalidParam
	default:
		return nil
	}
}

type FeeRatioRL struct {
	RentAmount         int64           `json:"rentAmount"`
	StakePerTrx        int64           `json:"stakePerTrx"`
	LiquidateThreshold decimal.Decimal `json:"liquidateThreshold"`
	RentalRate         decimal.Decimal `json:"rentalRate"`
	FeeRatio           decimal.Decimal `json:"feeRatio"`
	MinFee             decimal.Decimal `json:"minFee"`
	CurFeeRatio        decimal.Decimal `json:"curFeeRatio"`
	RentFee            decimal.Decimal `json:"rentFee"`
	PrePayFee          float64         `json:"prePayFee"`
}

var (
	_ internal.Conformer = (*FeeRatioMeta)(nil)
)

type FeeRatioService interface {
	// FeeRatio calculates the fee ratio based on the provided metadata
	FeeRatio(ctx context.Context, req *FeeRatioMeta) (*FeeRatioRL, error)
}
