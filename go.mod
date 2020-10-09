module github.com/chainapsis/cosmos-sdk-interchain-account

go 1.14

require (
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200921130040-27db2cf89772
	github.com/ghodss/yaml v1.0.0
	github.com/gibson042/canonicaljson-go v1.0.3 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.8
	github.com/regen-network/cosmos-proto v0.3.0
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tendermint v0.34.0-rc3.0.20200907055413-3359e0bf2f84
	github.com/tendermint/tm-db v0.6.2
	google.golang.org/grpc v1.32.0
	gopkg.in/yaml.v2 v2.3.0
	rsc.io/quote/v3 v3.1.0 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
