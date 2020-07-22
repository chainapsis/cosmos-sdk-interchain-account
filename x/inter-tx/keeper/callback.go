package keeper

import (
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/inter-tx/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (keeper Keeper) OnAccountCreated(ctx sdk.Context, sourcePort, sourceChannel string, address sdk.AccAddress) {
	receiver := keeper.popAddressFromRegistrationQueue(ctx, sourcePort, sourceChannel)

	if !receiver.Empty() {
		store := ctx.KVStore(keeper.storeKey)

		key := types.KeyRegisteredAccount(sourcePort, sourceChannel, receiver)
		store.Set(key, address)
	}
}
