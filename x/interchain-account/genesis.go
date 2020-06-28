package interchain_account

import (
	"fmt"

	"github.com/chainapsis/cosmos-sdk-interchain-account/x/interchain-account/keeper"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/interchain-account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, state types.GenesisState) {
	if !keeper.IsBound(ctx, state.PortID) {
		err := keeper.BindPort(ctx, state.PortID)
		if err != nil {
			panic(fmt.Sprintf("could not claim port capability: %v", err))
		}
	}
}

// ExportGenesis exports transfer module's portID into its geneis state
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	portID := keeper.GetPort(ctx)

	return types.GenesisState{
		PortID: portID,
	}
}
