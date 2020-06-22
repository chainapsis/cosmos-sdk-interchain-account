package inter_tx

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/chainapsis/cosmos-sdk-interchain-account/x/inter-tx/client/cli"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/inter-tx/client/rest"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/inter-tx/keeper"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct {
	cdc codec.Marshaler
}

func (AppModuleBasic) Name() string {
	return ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return nil
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, bz json.RawMessage) error {
	return nil
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

type AppModule struct {
	AppModuleBasic

	keeper Keeper
}

func NewAppModule(cdc codec.Marshaler, keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},

		keeper: keeper,
	}
}

func (AppModule) Name() string {
	return ModuleName
}

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// noop
}

// Route returns the message routing key for the staking module.
func (AppModule) Route() string {
	return RouterKey
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.keeper)
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	return nil
}

func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	// noop
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	// noop
	return []abci.ValidatorUpdate{}
}
