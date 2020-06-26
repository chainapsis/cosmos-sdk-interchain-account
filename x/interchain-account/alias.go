package interchain_account

import (
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/interchain-account/keeper"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/interchain-account/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	RouterKey    = types.RouterKey
	QuerierRoute = types.QuerierRoute
)

var (
	RegisterCodec = types.RegisterCodec
	NewKeeper     = keeper.NewKeeper
)

type (
	Keeper                                  = keeper.Keeper
	CounterpartyInfo                        = keeper.CounterpartyInfo
	InterchainAccountPacket                 = types.InterchainAccountPacket
	RegisterIBCAccountPacketData            = types.RegisterIBCAccountPacketData
	RunTxPacketData                         = types.RunTxPacketData
	RegisterIBCAccountPacketAcknowledgement = types.RegisterIBCAccountPacketAcknowledgement
	RunTxPacketAcknowledgement              = types.RunTxPacketAcknowledgement
)
