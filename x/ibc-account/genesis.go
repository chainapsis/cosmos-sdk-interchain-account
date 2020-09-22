package ibc_account

import (
	"fmt"

	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/keeper"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, state types.GenesisState) {
	if !keeper.IsBound(ctx, state.PortId) {
		err := keeper.BindPort(ctx, state.PortId)
		if err != nil {
			panic(fmt.Sprintf("could not claim port capability: %v", err))
		}
	}
}

// ExportGenesis exports transfer module's portID into its geneis state
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	portID := keeper.GetPort(ctx)

	return &types.GenesisState{
		PortId: portID,
	}
}
