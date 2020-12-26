package mock

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/testing/mock/client/cli"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/testing/mock/keeper"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/testing/mock/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is the mock AppModuleBasic.
type AppModuleBasic struct{}

// Name implements AppModuleBasic interface.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec implements AppModuleBasic interface.
func (AppModuleBasic) RegisterLegacyAminoCodec(*codec.LegacyAmino) {}

// RegisterInterfaces implements AppModuleBasic interface.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis implements AppModuleBasic interface.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return nil
}

// ValidateGenesis implements the AppModuleBasic interface.
func (AppModuleBasic) ValidateGenesis(codec.JSONMarshaler, client.TxEncodingConfig, json.RawMessage) error {
	return nil
}

// RegisterRESTRoutes implements AppModuleBasic interface.
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the bank module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// noop
}

// GetTxCmd implements AppModuleBasic interface.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd implements AppModuleBasic interface.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

// AppModule represents the AppModule for the mock module.
type AppModule struct {
	AppModuleBasic
	mockKeeper keeper.Keeper
}

// NewAppModule returns a mock AppModule instance.
func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{
		mockKeeper: k,
	}
}

// RegisterInvariants implements the AppModule interface.
func (AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route implements the AppModule interface.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.ModuleName, NewHandler(am.mockKeeper))
}

// QuerierRoute implements the AppModule interface.
func (AppModule) QuerierRoute() string {
	return ""
}

// LegacyQuerierHandler implements the AppModule interface.
func (am AppModule) LegacyQuerierHandler(*codec.LegacyAmino) sdk.Querier {
	return nil
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// noop
}

// InitGenesis implements the AppModule interface.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis implements the AppModule interface.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	return nil
}

// BeginBlock implements the AppModule interface
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
}

// EndBlock implements the AppModule interface
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
