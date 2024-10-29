package repos

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"justlend/internal"
	"justlend/internal/derrors"
	"justlend/internal/justlend"
	"justlend/internal/tron"
	"math/big"
)

func (ls *Service) RentResource(ctx context.Context,
	req *justlend.RentResourceMeta) (_ *justlend.RentResourceRL, err error) {
	defer derrors.WrapStack(&err, "ls.RentResource()")

	fee, err := ls.FeeRatio(ctx, &justlend.FeeRatioMeta{
		PrivateKey: req.PrivateKey,
		Type:       req.Type,
		Energy:     req.Amount,
	})

	stakePerTrx := tron.ToSUN(float64(fee.StakePerTrx))

	data := []byte{}
	methodId, _ := hexutil.Decode(justlend.RentResourceABI)
	paddedReceive := common.LeftPadBytes(internal.DecodeCheck(req.Receive)[1:], 32)
	paddedAmount := common.LeftPadBytes(new(big.Int).SetInt64(stakePerTrx).Bytes(), 32)
	paddedResourceType := common.LeftPadBytes(new(big.Int).SetInt64(int64(req.Type)).Bytes(), 32)
	data = append(data, methodId...)
	data = append(data, paddedReceive...)
	data = append(data, paddedAmount...)
	data = append(data, paddedResourceType...)
	result, txId, err := ls.tron.TriggerConstantContract(
		ctx,
		justlend.JustLendContract,
		data,
		req.PrivateKey,
		tron.ToSUN(fee.PrePayFee),
	)
	if err != nil {
		return nil, err
	}
	if _, err = ls.tron.BroadcastTransaction(ctx, result.Transaction); err != nil {
		return nil, err
	}
	return &justlend.RentResourceRL{
		TxId:        txId,
		StakePerTrx: stakePerTrx,
	}, nil
}
