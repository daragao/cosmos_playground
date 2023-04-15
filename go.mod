module github.com/daragao/cosmos_playground

go 1.19

require (
	cosmossdk.io/api v0.3.1
	google.golang.org/grpc v1.54.0
	google.golang.org/protobuf v1.29.1
)

require (
	github.com/cosmos/cosmos-proto v1.0.0-beta.2 // indirect
	github.com/cosmos/gogoproto v1.4.6 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230216225411-c8e22ba71e44 // indirect
)

replace (
	// osmosis-patched wasmd
	// ToDo: replace the commit hash with v0.30.0-osmo-v14 once the version is tagged
	github.com/CosmWasm/wasmd => github.com/osmosis-labs/wasmd v0.30.0-osmo-v15
	// dragonberry
	github.com/confio/ics23/go => github.com/cosmos/cosmos-sdk/ics23/go v0.8.0
	// Our cosmos-sdk branch is:  https://github.com/osmosis-labs/cosmos-sdk, current branch: v15.x. Direct commit link: https://github.com/osmosis-labs/cosmos-sdk/commit/44b40d47f3108c29f07fd115e5a92b387fb7a6bd
	github.com/cosmos/cosmos-sdk => github.com/osmosis-labs/cosmos-sdk v0.45.1-0.20230228211301-44b40d47f310
	// use cosmos-compatible protobufs
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	// Informal Tendermint fork
	github.com/tendermint/tendermint => github.com/informalsystems/tendermint v0.34.24
	// use grpc compatible with cosmos protobufs
	google.golang.org/grpc => google.golang.org/grpc v1.33.2

)
