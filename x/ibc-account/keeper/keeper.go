package keeper

import (
	"fmt"

	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
	ibcexported "github.com/cosmos/cosmos-sdk/x/ibc/exported"

	codec "github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/types"
)

const (
	// DefaultPacketTimeoutHeight is the default packet timeout height relative
	// to the current block height. The timeout is disabled when set to 0.
	DefaultPacketTimeoutHeight = 1000 // NOTE: in blocks

	// DefaultPacketTimeoutTimestamp is the default packet timeout timestamp relative
	// to the current block timestamp. The timeout is disabled when set to 0.
	DefaultPacketTimeoutTimestamp = 0 // NOTE: in nanoseconds
)

func SerializeCosmosTx(cdc codec.BinaryMarshaler, registry codectypes.InterfaceRegistry) func(data interface{}) ([]byte, error) {
	return func(data interface{}) ([]byte, error) {
		msgs := make([]sdk.Msg, 0)
		switch data := data.(type) {
		case sdk.Msg:
			msgs = append(msgs, data)
		case []sdk.Msg:
			msgs = append(msgs, data...)
		default:
			return nil, types.ErrInvalidOutgoingData
		}

		msgAnys := make([]*codectypes.Any, len(msgs))

		for i, msg := range msgs {
			var err error
			msgAnys[i], err = codectypes.NewAnyWithValue(msg)
			if err != nil {
				return nil, err
			}
		}

		txBody := &types.IBCTxBody{
			Messages: msgAnys,
		}

		txRaw := &types.IBCTxRaw{
			BodyBytes: cdc.MustMarshalBinaryBare(txBody),
		}

		bz, err := cdc.MarshalBinaryBare(txRaw)
		if err != nil {
			return nil, err
		}

		return bz, nil
	}
}

type CounterpartyInfo struct {
	// This method used to marshal transaction for counterparty chain.
	SerializeTx func(data interface{}) ([]byte, error)
}

// Keeper defines the IBC transfer keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.BinaryMarshaler

	// Key can be chain type which means what blockchain framework the host chain was built on or just direct chain id.
	counterpartyInfos map[string]CounterpartyInfo

	hook types.IBCAccountHooks

	channelKeeper types.ChannelKeeper
	portKeeper    types.PortKeeper
	accountKeeper types.AccountKeeper

	scopedKeeper capabilitykeeper.ScopedKeeper

	router types.Router
}

// NewKeeper creates a new IBC account Keeper instance
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey,
	counterpartyInfos map[string]CounterpartyInfo, channelKeeper types.ChannelKeeper, portKeeper types.PortKeeper,
	accountKeeper types.AccountKeeper, scopedKeeper capabilitykeeper.ScopedKeeper, router types.Router,
) Keeper {
	return Keeper{
		storeKey:          key,
		cdc:               cdc,
		counterpartyInfos: counterpartyInfos,
		channelKeeper:     channelKeeper,
		portKeeper:        portKeeper,
		accountKeeper:     accountKeeper,
		scopedKeeper:      scopedKeeper,
		router:            router,
	}
}

func (k Keeper) AddCounterpartyInfo(typ string, info CounterpartyInfo) {
	k.counterpartyInfos[typ] = info
}

func (k Keeper) GetCounterpartyInfo(typ string) (CounterpartyInfo, bool) {
	info, ok := k.counterpartyInfos[typ]
	return info, ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s-%s", host.ModuleName, types.ModuleName))
}

// IsBound checks if the interchain account module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	// Set the portID into our store so we can retrieve it later
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.PortKey), []byte(portID))

	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// GetPort returns the portID for the ibc account module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get([]byte(types.PortKey)))
}

// ClaimCapability allows the transfer module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}
