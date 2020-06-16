package interchain_account

import (
	"encoding/json"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/capability"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/chainapsis/cosmos-sdk-interchain-account/x/interchain-account/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	porttypes "github.com/cosmos/cosmos-sdk/x/ibc/05-port/types"
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
	return ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
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
	keeper Keeper
}

func NewAppModule(k Keeper) AppModule {
	return AppModule{
		keeper: k,
	}
}

func (AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// noop
}

func (AppModule) Route() string {
	return RouterKey
}

func (AppModule) NewHandler() sdk.Handler {
	return nil
}

func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

func (AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

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
	chanCap *capability.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) error {
	if order != ibctypes.ORDERED {
		return types.ErrInvalidOrder
	}

	boundPort := am.keeper.GetPort(ctx)
	if boundPort != portID {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	}

	// TODO: Check version
	// if version != types.Version {
	// 	return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, "ics20-1")
	// }

	if err := am.keeper.ClaimCapability(ctx, chanCap, ibctypes.ChannelCapabilityPath(portID, channelID)); err != nil {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error())
	}

	return nil
}

func (am AppModule) OnChanOpenTry(
	ctx sdk.Context,
	order ibctypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capability.Capability,
	counterparty channeltypes.Counterparty,
	version,
	counterpartyVersion string,
) error {
	if order != ibctypes.ORDERED {
		return types.ErrInvalidOrder
	}

	boundPort := am.keeper.GetPort(ctx)
	if boundPort != portID {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected %s", portID, boundPort)
	}

	// TODO: Check version
	// if version != types.Version {
	// 	return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid version: %s, expected %s", version, "ics20-1")
	// }

	// TODO: Check counterparty version
	// if counterpartyVersion != types.Version {
	// 	return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid counterparty version: %s, expected %s", counterpartyVersion, "ics20-1")
	// }

	// Claim channel capability passed back by IBC module
	if err := am.keeper.ClaimCapability(ctx, chanCap, ibctypes.ChannelCapabilityPath(portID, channelID)); err != nil {
		return sdkerrors.Wrap(channel.ErrChannelCapabilityNotFound, err.Error())
	}

	// TODO: escrow
	return nil
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
	var data InterchainAccountPacket
	// TODO: Remove the usage of global variable "ModuleCdc"
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal interchain account packet data: %s", err.Error())
	}

	err := am.keeper.OnRecvPacket(ctx, packet, data)

	switch data.(type) {
	case RegisterIBCAccountPacketData:
		acknowledgement := RegisterIBCAccountPacketAcknowledgement{
			Success: err == nil,
		}

		if err := am.keeper.PacketExecuted(ctx, packet, acknowledgement.GetBytes()); err != nil {
			return nil, err
		}
		return &sdk.Result{
			Events: ctx.EventManager().Events().ToABCIEvents(),
		}, nil
	case RunTxPacketData:
		acknowledgement := RunTxPacketAcknowledgement{
			Code: 0,
		}
		if err != nil {
			acknowledgement = RunTxPacketAcknowledgement{
				Code: 1,
			}
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
	// TODO
	return nil, nil
}

func (am AppModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) (*sdk.Result, error) {
	// TODO
	return nil, nil
}
