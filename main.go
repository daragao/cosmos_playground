package main

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/unknownproto"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"

	"context"

	"google.golang.org/grpc"

	tendermint "cosmossdk.io/api/cosmos/base/tendermint/v1beta1"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	txType "github.com/cosmos/cosmos-sdk/types/tx"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
)

func main() {
	url := "grpc.osmosis.zone:9090"
	// Create a connection to the gRPC server.
	grpcConn, err := grpc.Dial(
		url,                 // your gRPC server address.
		grpc.WithInsecure(), // The Cosmos SDK doesn't support any transport security mechanism.
		// This instantiates a general gRPC codec which handles proto bytes. We pass in a nil interface registry
		// if the request/response types contain interface instead of 'nil' you should pass the application specific codec.
		//grpc.WithDefaultCallOptions(grpc.ForceCodec(codec.NewProtoCodec(nil).GRPCCodec())),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer grpcConn.Close()

	cli := tendermint.NewServiceClient(grpcConn)

	lastBlock, err := cli.GetLatestBlock(
		context.Background(),
		&tendermint.GetLatestBlockRequest{},
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	block := lastBlock.GetBlock()
	blockData := block.GetData()
	blockTxs := blockData.GetTxs()
	fmt.Printf("Block: %3.d Total txs: %d\n", block.GetHeader().GetHeight(), len(blockTxs))

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	defaultDecodeTx := authTx.DefaultTxDecoder(cdc)

	for i, txB := range blockTxs {
		//_, err := decodeTx(cdc, txB)
		tx, err := defaultDecodeTx(txB)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("========\n%d\n%#v\n========\n", i, tx)
		return
	}
}

func decodeTx(cdc codec.ProtoCodecMarshaler, txBytes []byte) (*txType.TxRaw, error) {
	var raw tx.TxRaw

	// reject all unknown proto fields in the root TxRaw
	err := unknownproto.RejectUnknownFieldsStrict(txBytes, &raw, cdc.InterfaceRegistry())
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}

	err = cdc.Unmarshal(txBytes, &raw)
	if err != nil {
		return nil, err
	}

	var any codectypes.Any
	err = cdc.Unmarshal(raw.GetBodyBytes(), &any)
	if err != nil {
		return nil, err
	}
	fmt.Printf("---\n%#v\n---", any)
	return &raw, nil
}
