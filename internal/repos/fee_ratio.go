package repos

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	"justlend/internal"
	"justlend/internal/derrors"
	"justlend/internal/justlend"
	"justlend/internal/protos/core"
	"math/big"
)

func (ls *Service) FeeRatio(ctx context.Context, req *justlend.FeeRatioMeta) (_ *justlend.FeeRatioRL, err error) {
	defer derrors.WrapStack(&err, "ls.FeeRatio()")

	pub := internal.PrivateKeyToPublicKey(req.PrivateKey)
	address := internal.EncodeCheck(internal.PublicKeyToTronAddress(pub))

	_liquidateThreshold, _ := ls.liquidateThreshold(ctx, address, req.PrivateKey)

	stakePerTrx, err := ls.tron.CalStackEnergy(ctx, address, req.Energy, false)
	if err != nil {
		return nil, err
	}

	rentalRate, _ := ls.getRentalRate(ctx, stakePerTrx, req.PrivateKey, req.Type)
	_feeRatio, _ := ls.feeRatio(ctx, address, req.PrivateKey)
	_minFee, _ := ls.minFee(ctx, address, req.PrivateKey)

	curFeeRatio := _feeRatio.Mul(decimal.NewFromInt(stakePerTrx))
	feeRatio := decimal.Zero
	if _minFee.Compare(curFeeRatio) > 0 {
		feeRatio = _minFee
	} else {
		feeRatio = curFeeRatio
	}
	countFeeDuration := decimal.NewFromInt(172800)
	rentFee := decimal.NewFromInt(stakePerTrx).
		Mul(rentalRate).
		Mul(countFeeDuration).
		Add(_liquidateThreshold)
	prePayFee, _ := rentFee.Add(feeRatio).Float64()
	return &justlend.FeeRatioRL{
		RentAmount:         req.Energy,
		StakePerTrx:        stakePerTrx,
		LiquidateThreshold: _liquidateThreshold,
		RentalRate:         rentalRate,
		FeeRatio:           feeRatio,
		MinFee:             _minFee,
		CurFeeRatio:        curFeeRatio,
		RentFee:            rentFee,
		PrePayFee:          prePayFee,
	}, nil
}

func (ls *Service) getRentalRate(ctx context.Context,
	value int64,
	privateKey string,
	rt core.ResourceCode) (decimal.Decimal, error) {
	data := []byte{}
	methodId, _ := hexutil.Decode(justlend.GetRentalRateABI)
	resourceType := common.LeftPadBytes(new(big.Int).SetInt64(int64(rt)).Bytes(), 32)
	amount := common.LeftPadBytes(new(big.Int).SetInt64(value).Bytes(), 32)
	data = append(data, methodId...)
	data = append(data, amount...)
	data = append(data, resourceType...)

	result, _, err := ls.tron.TriggerConstantContract(
		ctx,
		justlend.JustLendContract,
		data,
		privateKey,
		0,
	)
	if err != nil {
		return decimal.Zero, nil
	}
	bigInt := new(big.Int).SetBytes(result.GetConstantResult()[0]).Int64()
	d := decimal.NewFromInt(bigInt)
	rentalRate := d.Div(decimal.NewFromInt(justlend.TokenDefaultPrecision))
	return rentalRate, nil
}

func (ls *Service) liquidateThreshold(ctx context.Context, address, privateKey string) (decimal.Decimal, error) {

	data := []byte{}

	methodId, _ := hexutil.Decode(justlend.LiquidateThresholdABI)

	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(internal.DecodeCheck(address)[1:], 32)...)
	data = append(data, common.LeftPadBytes(internal.DecodeCheck(justlend.JustLendContract)[1:], 32)...)
	result, _, err := ls.tron.TriggerConstantContract(ctx, justlend.JustLendContract, data, privateKey, 0)
	if err != nil {
		return decimal.Zero, err
	}
	newLiquidateThreshold := new(big.Int).SetBytes(result.GetConstantResult()[0]).Int64()
	return decimal.NewFromInt(newLiquidateThreshold).Div(decimal.NewFromInt(1000000)), nil
}

func (ls *Service) minFee(ctx context.Context, address, privateKey string) (decimal.Decimal, error) {
	data := []byte{}
	methodId, _ := hexutil.Decode(justlend.MinFeeABI)
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(internal.DecodeCheck(address)[1:], 32)...)
	data = append(data, common.LeftPadBytes(internal.DecodeCheck(justlend.JustLendContract)[1:], 32)...)
	result, _, err := ls.tron.TriggerConstantContract(ctx, justlend.JustLendContract, data, privateKey, 0)
	if err != nil {
		return decimal.Zero, err
	}
	newLiquidateThreshold := new(big.Int).SetBytes(result.GetConstantResult()[0]).Int64()
	return decimal.NewFromInt(newLiquidateThreshold).Div(decimal.NewFromInt(1000000)), nil
}

func (ls *Service) feeRatio(ctx context.Context, address, privateKey string) (decimal.Decimal, error) {
	data := []byte{}
	methodId, _ := hexutil.Decode(justlend.FeeRatioABI)
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(internal.DecodeCheck(address)[1:], 32)...)
	data = append(data, common.LeftPadBytes(internal.DecodeCheck(justlend.JustLendContract)[1:], 32)...)
	result, _, err := ls.tron.TriggerConstantContract(ctx, justlend.JustLendContract, data, privateKey, 0)
	if err != nil {
		return decimal.Zero, err
	}
	feeRatioNotPrecision := new(big.Int).SetBytes(result.GetConstantResult()[0]).Int64()
	return decimal.NewFromInt(feeRatioNotPrecision).Div(decimal.NewFromUint64(justlend.TokenDefaultPrecision)), nil
}
