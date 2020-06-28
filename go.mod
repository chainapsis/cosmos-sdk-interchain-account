module github.com/chainapsis/cosmos-sdk-interchain-account

go 1.14

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200626163604-43837b16e768
	github.com/gogo/protobuf v1.3.1
	github.com/gorilla/mux v1.7.4
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tendermint v0.33.5
	github.com/tendermint/tm-db v0.5.1
	github.com/regen-network/cosmos-proto v0.3.0
	google.golang.org/grpc v1.30.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
