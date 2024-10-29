package tron

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"justlend/internal"
	"justlend/internal/config"
	"justlend/internal/derrors"
	"justlend/internal/protos/api"
	"justlend/internal/protos/core"
	"log"
	"math"
)

var (
	defaultTrongRPCEndpoint = config.GetEnv("TRON_GRPC_ENDPOINT", "34.220.77.106:50051")
)

type Endpoint struct {
	node   string           // endpoint of the Tron node
	grpc   *grpc.ClientConn // client connection to the Wallet service
	wallet api.WalletClient // client API for wallet service
}

// NewEndpoint creates a new endpoint to the Tron node.
func NewEndpoint() (*Endpoint, error) {
	// Create a new gRPC client connection to the Tron node.
	conn, err := grpc.NewClient(
		defaultTrongRPCEndpoint,
		// Use insecure credentials for now.
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		log.Panic(err)
	}
	return &Endpoint{
		node:   defaultTrongRPCEndpoint,
		grpc:   conn,
		wallet: api.NewWalletClient(conn),
	}, nil
}

const SUNPerTRX = 1000000

func ToSUN(trx float64) int64 {
	return int64(trx * SUNPerTRX)
}

// GetAccountResource Retrieve the account resource using the wallet client
func (e *Endpoint) GetAccountResource(ctx context.Context, address string) (*api.AccountResourceMessage, error) {
	return e.wallet.GetAccountResource(ctx, &core.Account{Address: internal.DecodeCheck(address)})
}

func (e *Endpoint) TriggerConstantContract(ctx context.Context,
	contract string,
	data []byte,
	privateKey string,
	callValue int64) (*api.TransactionExtention, string, error) {

	pub := internal.PrivateKeyToPublicKey(privateKey)
	address := internal.EncodeCheck(internal.PublicKeyToTronAddress(pub))

	transferContract := new(core.TriggerSmartContract)
	transferContract.OwnerAddress = internal.DecodeCheck(address)
	transferContract.ContractAddress = internal.DecodeCheck(contract)
	transferContract.Data = data
	transferContract.CallValue = callValue

	transferTransactionEx, err := e.wallet.TriggerConstantContract(ctx, transferContract)
	if err != nil {
		return nil, "", err
	}
	transferTransaction := transferTransactionEx.Transaction
	if transferTransaction == nil ||
		len(transferTransaction.GetRawData().GetContract()) == 0 {
		return nil, "", fmt.Errorf("transfer error: invalid transaction")
	}
	var txId []byte
	transferTransaction.RawData.FeeLimit = 200000000
	txId, err = signTransaction(transferTransaction, privateKey)
	if err != nil {
		return nil, "", err
	}
	return transferTransactionEx, hex.EncodeToString(txId), nil
}

func (e *Endpoint) CalStackEnergy(ctx context.Context, owner string, energy int64, toSUN bool) (int64, error) {
	return e.StackEnergy(ctx, owner, energy, toSUN)
}

// StackEnergy is a method that calculates the obtained energy based on the given energy value and the account's resource information.
// It retrieves the account's resource and performs calculations based on the energy weight and energy limit.
// The function returns the obtained energy value if successful, otherwise an error is returned.
func (e *Endpoint) StackEnergy(ctx context.Context,
	owner string,
	energy int64,
	toSUN bool,
) (int64, error) {
	// Retrieve the account's resource information
	resource, err := e.GetAccountResource(ctx, owner)
	if err != nil || resource == nil {
		return -1, derrors.Forbidden
	} else if energyWeight := decimal.NewFromInt(resource.GetTotalEnergyWeight()); energyWeight.IsZero() {
		// Check for invalid energy weight
		return -1, fmt.Errorf(`stackEnergy: invalid energy weight(%v)`, energyWeight)
	} else if energyLimit := decimal.NewFromInt(resource.GetTotalEnergyLimit()); energyLimit.IsZero() {
		// Check for invalid energy limit
		return -1, fmt.Errorf(`stackEnergy: invalid energy limit(%v)`, energyLimit)
	} else {
		// Calculate the obtained energy based on the energy value, energy limit, and energy weight
		proportion, _ := energyLimit.Div(energyWeight).Float64()
		front := math.Ceil(float64(energy) / proportion)
		if toSUN {
			return int64(math.Ceil(front * 1000000)), nil
		} else {
			return int64(math.Ceil(front)), nil
		}
	}
}

// BroadcastTransaction is a method that broadcasts a transaction to the wallet.
// It calls the BroadcastTransaction method of the wallet to perform the broadcasting process.
func (e *Endpoint) BroadcastTransaction(ctx context.Context, transaction *core.Transaction) (bool, error) {
	if reply, err := e.wallet.BroadcastTransaction(ctx, transaction); err != nil {
		return false, err
	} else if !reply.Result {
		return false, fmt.Errorf("%s", reply.String())
	}
	return true, nil
}

func (e *Endpoint) Close() error {
	return e.grpc.Close()
}
