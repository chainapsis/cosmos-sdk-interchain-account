package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (registerData RegisterIBCAccountPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(registerData))
}

func (runTxData RunTxPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(runTxData))
}

func (registerAck RegisterIBCAccountPacketAcknowledgement) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(registerAck))
}

func (runtxAck RunTxPacketAcknowledgement) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(runtxAck))
}
