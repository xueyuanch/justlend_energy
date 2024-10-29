package repos

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"justlend/internal"
	"justlend/internal/derrors"
	"justlend/internal/justlend"
	"math/big"
)

func (ls *Service) ReturnResource(ctx context.Context,
	req *justlend.ReturnResourceMeta) (_ *justlend.ReturnResourceRL, err error) {

	defer derrors.WrapStack(&err, "ls.ReturnResource()")

	data := []byte{}
	methodId, _ := hexutil.Decode(justlend.ReturnResourceABI)
	data = append(data, methodId...)
	data = append(data, common.LeftPadBytes(internal.DecodeCheck(req.Receive)[1:], 32)...)
	data = append(data, common.LeftPadBytes(new(big.Int).SetInt64(req.StakePerTrx).Bytes(), 32)...)
	data = append(data, common.LeftPadBytes(new(big.Int).SetInt64(int64(req.Type)).Bytes(), 32)...)
	result, txId, err := ls.tron.TriggerConstantContract(
		ctx,
		justlend.JustLendContract,
		data,
		req.PrivateKey,
		0,
	)
	if err != nil {
		return nil, err
	}
	if _, err = ls.tron.BroadcastTransaction(ctx, result.Transaction); err != nil {
		return nil, err
	}
	return &justlend.ReturnResourceRL{
		TxId: txId,
	}, nil
}
