package keeper

import (
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
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

func (keeper Keeper) TryRunTxMsgSend(ctx sdk.Context, sourcePort, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, fromAddr, toAddr sdk.AccAddress, amount sdk.Coins) error {
	msg := banktypes.NewMsgSend(fromAddr, toAddr, amount)
	_, err := keeper.ibcAccountKeeper.TryRunTx(ctx, sourcePort, sourceChannel, "cosmos-sdk", msg, timeoutHeight, timeoutTimestamp)
	return err
}
