package keeper

import (
	"fmt"

	"github.com/chainapsis/cosmos-sdk-interchain-account/x/inter-tx/types"
	ia "github.com/chainapsis/cosmos-sdk-interchain-account/x/interchain-account"
	iatypes "github.com/chainapsis/cosmos-sdk-interchain-account/x/interchain-account/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc               codec.Marshaler
	counterpartyTxCdc *codec.Codec
	storeKey          sdk.StoreKey

	iaKeeper ia.Keeper
}

func NewKeeper(cdc codec.Marshaler, counterpartyTxCdc *codec.Codec, storeKey sdk.StoreKey, iaKeeper ia.Keeper) Keeper {
	return Keeper{
		cdc:               cdc,
		counterpartyTxCdc: counterpartyTxCdc,
		storeKey:          storeKey,

		iaKeeper: iaKeeper,
	}
}

//nolint:interfacer
func (keeper Keeper) RegisterInterchainAccount(ctx sdk.Context, sender sdk.AccAddress, sourcePort, sourceChannel string) error {
	salt := keeper.GetIncrementalSalt(ctx)
	err := keeper.iaKeeper.CreateInterchainAccount(ctx, sourcePort, sourceChannel, salt)
	if err != nil {
		return err
	}

	address := keeper.iaKeeper.GenerateAddress(iatypes.GetIdentifier(sourcePort, sourceChannel), salt)

	kvStore := ctx.KVStore(keeper.storeKey)
	prefixStore := prefix.NewStore(kvStore, []byte("ia/"))

	key := []byte(fmt.Sprintf("%s/%s/%s", sourcePort, sourceChannel, sender.String()))
	if prefixStore.Has(key) {
		return types.ErrIAAccountAlreadyExist
	}
	prefixStore.Set(key, address)

	ctx.EventManager().EmitEvent(sdk.NewEvent("register-interchain-account",
		sdk.NewAttribute("expected-address", sdk.AccAddress(address).String()),
		sdk.NewAttribute("salt", salt)))

	return nil
}

func (keeper Keeper) RunMsg(ctx sdk.Context, msgBytes []byte, sender sdk.AccAddress, sourcePort, sourceChannel string) error {
	// TODO: Use counterpart chain's codec.
	var msg sdk.Msg
	err := keeper.counterpartyTxCdc.UnmarshalBinaryBare(msgBytes, &msg)
	if err != nil {
		return err
	}

	return keeper.RunTx(ctx, []sdk.Msg{msg}, sender, sourcePort, sourceChannel)
}

func (keeper Keeper) RunTx(ctx sdk.Context, msgs []sdk.Msg, sender sdk.AccAddress, sourcePort, sourceChannel string) error {
	_, err := keeper.GetInterchainAccount(ctx, sender, sourcePort, sourceChannel)
	if err != nil {
		return err
	}

	err = keeper.iaKeeper.RequestRunTx(ctx, sourcePort, sourceChannel, iatypes.CosmosSdkChainType, msgs)
	if err != nil {
		return err
	}

	return nil
}

//nolint:interfacer
func (keeper Keeper) GetInterchainAccount(ctx sdk.Context, address sdk.AccAddress, sourcePort, sourceChannel string) ([]byte, error) {
	kvStore := ctx.KVStore(keeper.storeKey)
	prefixStore := prefix.NewStore(kvStore, []byte("ia/"))

	key := []byte(fmt.Sprintf("%s/%s/%s", sourcePort, sourceChannel, address.String()))
	if !prefixStore.Has(key) {
		return []byte{}, types.ErrIAAccountNotExist
	}
	bz := prefixStore.Get(key)

	return bz, nil
}

func (keeper Keeper) GetIncrementalSalt(ctx sdk.Context) string {
	kvStore := ctx.KVStore(keeper.storeKey)

	key := []byte("salt")

	salt := 0
	if kvStore.Has(key) {
		keeper.cdc.MustUnmarshalJSON(kvStore.Get(key), &salt)
		salt++
	}

	bz := keeper.cdc.MustMarshalJSON(salt)
	kvStore.Set(key, bz)

	return string(bz)
}
