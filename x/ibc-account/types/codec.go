package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	commitmenttypes "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment/types"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(IBCAccountPacketData{}, "ibcaccount/IBCAccountPacketData", nil)
	cdc.RegisterConcrete(IBCAccountPacketAcknowledgement{}, "ibcaccount/IBCAccountPacketAcknowledgement", nil)
}

var (
	amino = codec.New()

	ModuleCdc = codec.NewHybridCodec(amino, cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(amino)
	channel.RegisterCodec(amino)
	commitmenttypes.RegisterCodec(amino)
	amino.Seal()
}
