package main

import (
	"fmt"
	"os"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/unknownproto"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"

	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protowire"

	tendermint "cosmossdk.io/api/cosmos/base/tendermint/v1beta1"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	txType "github.com/cosmos/cosmos-sdk/types/tx"
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

	/*
		registry := codectypes.NewInterfaceRegistry()
		cdc := codec.NewProtoCodec(registry)
		defaultDecodeTx := authTx.DefaultTxDecoder(cdc)
	*/

	for _, txB := range blockTxs {
		/*
			//_, err := decodeTx(cdc, txB)
			tx, err := defaultDecodeTx(txB)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("========\n%d\n%#v\n========\n", i, tx)
		*/

		unknownTxBytes(txB)
		return
	}
}

func writeBytes(filename string, b []byte) {
	err := os.WriteFile(filename, b, 0644)
	if err != nil {
		panic(err)
	}
}

func unknownTxBytes(txB []byte) {
	for len(txB) > 0 {
		fmt.Println("===========================")
		tagNum, wireType, tagLen := protowire.ConsumeTag(txB)
		fmt.Print("tagNum: ", tagNum)
		if wireType != 2 {
			panic(fmt.Errorf("unexpected wireType (%d)", wireType))
		}
		txB := txB[tagLen:]

		v, n := protowire.ConsumeVarint(txB)
		txB = txB[n:]
		fmt.Println("bytes length: ", v)

		fieldValue, n := protowire.ConsumeBytes(txB[:v])

		fmt.Println(v, n, string(fieldValue))
		unknownTxBytes(txB[v:])

		return
		/*
			any := new(codectypes.Any)
			if err := proto.Unmarshal(fieldValue, any); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(any)
		*/

		/*
			tagNum, wireType, fieldLen := protowire.ConsumeField(txB)
			if fieldLen < 0 {
				return
			}
			fieldB := txB[:fieldLen]
			txB = txB[fieldLen:]
			fmt.Println("-->", tagNum, wireType, fieldLen)

			tagNum, wireType, tagLen := protowire.ConsumeTag(fieldB)
			tagB := fieldB[:tagLen]
			fieldB = fieldB[tagLen:]
			fmt.Println("-->", tagNum, wireType, tagLen, tagB)

			switch wireType {
			case 0:
			case 1:
				i, o := protowire.ConsumeVarint(fieldB)
				varint := fieldB[:o]
				fieldB = fieldB[o:]
				fmt.Println("VARINT -->", i, o, varint)
				break
			case 2:
				v, _ := protowire.ConsumeBytes(fieldB)
				unknownTxBytes(v)
				//fmt.Println("-->", string(v), bLen)
			}
		*/

		/*
			unknownTxBytes(fieldB)
			any := new(codectypes.Any)
			if err := proto.Unmarshal(fieldB, any); err != nil {
				fmt.Println(err)
				return
			}
		*/

		//fmt.Println("-->", any.TypeUrl)
		fmt.Println("===========================")
		//fieldBytes = any.Value

		//n := protowire.ConsumeFieldValue(tagNum, wireType, txB[m:])
		//fmt.Println("-->", n, string(txB[m:m+n]))
		//fieldBytes := txB[m : m+n]
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
