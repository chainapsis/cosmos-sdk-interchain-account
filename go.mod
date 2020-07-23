module github.com/chainapsis/cosmos-sdk-interchain-account

go 1.14

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200722165112-e0a81d2cca8f
	github.com/gogo/protobuf v1.3.1
	github.com/gorilla/mux v1.7.4
	github.com/regen-network/cosmos-proto v0.3.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tendermint v0.33.6
	github.com/tendermint/tm-db v0.5.1
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.24.0 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
