package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type InterchainAccountPacket interface{}

const CosmosSdkChainType = "cosmos-sdk"

type InterchainAccountTx struct {
	Msgs []sdk.Msg
}
