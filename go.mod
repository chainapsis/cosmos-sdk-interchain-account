module github.com/chainapsis/cosmos-sdk-interchain-account

go 1.14

require (
	github.com/armon/go-metrics v0.3.3 // indirect
	github.com/confio/ics23-iavl v0.6.0 // indirect
	github.com/confio/ics23/go v0.0.0-20200604202538-6e2c36a74465 // indirect
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200530180557-ba70f4d4dc2e
	github.com/gogo/protobuf v1.3.1
	github.com/gorilla/mux v1.7.4
	github.com/otiai10/copy v1.2.0
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/regen-network/cosmos-proto v0.3.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/iavl v0.13.3 // indirect
	github.com/tendermint/tendermint v0.33.4
	github.com/tendermint/tm-db v0.5.1
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.24.0 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
