package ibc_account

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"

	"github.com/gogo/protobuf/grpc"

	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/keeper"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var (
	_ module.AppModule      = AppModule{}
	_ port.IBCModule        = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

func (AppModuleBasic) ValidateGenesis(_ codec.JSONMarshaler, _ json.RawMessage) error {
	return nil
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	// noop
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	// noop
	return nil
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// noop
	return nil
}

// RegisterInterfaceTypes registers module concrete types into protobuf Any.
func (AppModuleBasic) RegisterInterfaceTypes(registry cdctypes.InterfaceRegistry) {
	// noop
}

type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{
		keeper: k,
	}
}

func (AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// noop
}

func (AppModule) Route() string {
	return types.RouterKey
}

func (AppModule) NewHandler() sdk.Handler {
	return nil
}

func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

func (am AppModule) RegisterQueryService(grpc.Server) {}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// BeginBlock implements the AppModule interface
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {

}

// EndBlock implements the AppModule interface
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// Implement IBCModule callbacks
func (am AppModule) OnChanOpenInit(
	ctx sdk.Context,
	order ibctypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) error {
	return am.keeper.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, chanCap, counterparty, version)
}

func (am AppModule) OnChanOpenTry(
	ctx sdk.Context,
	order ibctypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version,
	counterpartyVersion string,
) error {
	return am.keeper.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, chanCap, counterparty, version, counterpartyVersion)
}

func (am AppModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyVersion string,
) error {
	// TODO
	// if counterpartyVersion != types.Version {
	//	return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid counterparty version: %s, expected %s", counterpartyVersion, "ics20-1")
	// }
	return nil
}

func (am AppModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

func (am AppModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Disallow user-initiated channel closing for interchain account channels
	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

func (am AppModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

func (am AppModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) (*sdk.Result, error) {
	var data types.IBCAccountPacketData
	// TODO: Remove the usage of global variable "ModuleCdc"
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal interchain account packet data: %s", err.Error())
	}

	err := am.keeper.OnRecvPacket(ctx, packet)

	switch data.Type {
	case types.Type_REGISTER:
		acknowledgement := types.IBCAccountPacketAcknowledgement{
			ChainID: ctx.ChainID(),
		}
		if err == nil {
			acknowledgement.Code = 0
		} else {
			acknowledgement.Code = 1
			acknowledgement.Error = err.Error()
		}

		if err := am.keeper.PacketExecuted(ctx, packet, acknowledgement.GetBytes()); err != nil {
			return nil, err
		}
		return &sdk.Result{
			Events: ctx.EventManager().Events().ToABCIEvents(),
		}, nil
	case types.Type_RUNTX:
		acknowledgement := types.IBCAccountPacketAcknowledgement{
			ChainID: ctx.ChainID(),
		}
		if err == nil {
			acknowledgement.Code = 0
		} else {
			acknowledgement.Code = 1
			acknowledgement.Error = err.Error()
		}

		if err := am.keeper.PacketExecuted(ctx, packet, acknowledgement.GetBytes()); err != nil {
			return nil, err
		}
		return &sdk.Result{
			Events: ctx.EventManager().Events().ToABCIEvents(),
		}, nil
	default:
		return nil, types.ErrUnknownPacketData
	}
}

func (am AppModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
) (*sdk.Result, error) {
	var ack types.IBCAccountPacketAcknowledgement
	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 interchain account packet acknowledgment: %v", err)
	}
	var data types.IBCAccountPacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 interchain account packet data: %s", err.Error())
	}

	if err := am.keeper.OnAcknowledgementPacket(ctx, packet, data, ack); err != nil {
		return nil, err
	}

	// TODO: Add events.

	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func (am AppModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) (*sdk.Result, error) {
	// TODO
	return nil, nil
}
