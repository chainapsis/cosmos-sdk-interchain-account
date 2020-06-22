package inter_tx

import (
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/inter-tx/keeper"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/inter-tx/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	RouterKey    = types.RouterKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper          = keeper.NewKeeper
	RegisterCodec      = types.RegisterCodec
	RegisterInterfaces = types.RegisterInterfaces
	NewMsgRegister     = types.NewMsgRegister
	NewMsgRunTx        = types.NewMsgRunTx
)

type (
	Keeper      = keeper.Keeper
	MsgRegister = types.MsgRegister
	MsgRunTx    = types.MsgRunTx
)
