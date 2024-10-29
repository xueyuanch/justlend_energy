package tron

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gogo/protobuf/proto"
	"justlend/internal/protos/core"
	"time"
)

func parsePrivateKey(private string) *ecdsa.PrivateKey {
	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		return nil
	}
	return privateKey
}

func signTransaction(transaction *core.Transaction, privateKey string) ([]byte, error) {
	transaction.GetRawData().Timestamp = time.Now().UnixNano() / 1000000
	rawData, err := proto.Marshal(transaction.GetRawData())
	if err != nil {
		return nil, err
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	contractList := transaction.GetRawData().GetContract()
	for range contractList {
		s, e := crypto.Sign(hash, parsePrivateKey(privateKey))
		if e != nil {
			return nil, e
		}
		transaction.Signature = append(transaction.Signature, s)
	}
	return hash, nil
}
