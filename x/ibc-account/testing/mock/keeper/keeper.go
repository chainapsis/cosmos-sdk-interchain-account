package keeper

import (
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/02-client/types"
)

type Keeper struct {
	ibcAccountKeeper keeper.Keeper
}

func NewKeeper(ibcAccountKeeper keeper.Keeper) Keeper {
	return Keeper{
		ibcAccountKeeper: ibcAccountKeeper,
	}
}

func (keeper Keeper) TryRegisterIBCAccount(ctx sdk.Context, sourcePort, sourceChannel string, salt []byte, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) error {
	return keeper.ibcAccountKeeper.TryRegisterIBCAccount(ctx, sourcePort, sourceChannel, salt, timeoutHeight, timeoutTimestamp)
}
