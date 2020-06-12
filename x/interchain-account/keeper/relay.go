package keeper

import (
	"bytes"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/interchain-account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"math"
)

func (k Keeper) RegisterIBCAccount(
	ctx sdk.Context,
	sourcePort,
	sourceChannel,
	salt string,
) error {
	identifier := types.GetIdentifier(sourcePort, sourceChannel)
	address, err := k.GenerateAddress(identifier, salt)
	if err != nil {
		return err
	}
	err = k.CreateAccount(ctx, address, identifier)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) CreateAccount(ctx sdk.Context, address sdk.AccAddress, identifier string) error {
	account := k.accountKeeper.GetAccount(ctx, address)
	// Don't block even if there is normal account,
	// because attackers can distrupt to create an interchain account
	// by sending some assets to estimated address in advance.
	if account != nil {
		if account.GetSequence() != 0 || account.GetPubKey() != nil {
			// If account is interchain account or is usable by someone.
			return sdkerrors.Wrap(types.ErrAccountAlreadyExist, account.String())
		}
		err := account.SetSequence(1)
		if err != nil {
			return err
		}
	} else {
		account = k.accountKeeper.NewAccountWithAddress(ctx, address)
		err := account.SetSequence(1)
		if err != nil {
			return err
		}
	}

	// Interchain accounts have the sequence "1" and nil public key.
	// Sequence never be increased without signing tx and sending this tx.
	// But, it is impossible to send tx without publishing the public key.
	// So, accounts that have the sequence "1" and nill public key are explicitly interchain accounts.
	k.accountKeeper.SetAccount(ctx, account)

	store := ctx.KVStore(k.storeKey)
	// Save the identifier for each address to check where the interchain account is made from.
	store.Set(address, []byte(identifier))

	return nil
}

// Determine account's address that will be created.
func (k Keeper) GenerateAddress(identifier string, salt string) ([]byte, error) {
	return tmhash.SumTruncated([]byte(identifier + salt)), nil
}

func (k Keeper) CreateInterchainAccount(ctx sdk.Context, sourcePort, sourceChannel, salt string) error {
	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrap(channeltypes.ErrChannelNotFound, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	channelCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return channel.ErrSequenceSendNotFound
	}

	packetData := types.RegisterIBCAccountPacketData{
		Salt: salt,
	}

	// TODO: Add timeout height and timestamp
	packet := channeltypes.NewPacket(
		packetData.GetBytes(),
		sequence,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		math.MaxUint64,
		0,
	)

	return k.channelKeeper.SendPacket(ctx, channelCap, packet)
}

func (k Keeper) RequestRunTx(ctx sdk.Context, sourcePort, sourceChannel, chainType string, data interface{}) error {
	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrap(channeltypes.ErrChannelNotFound, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	return k.CreateOutgoingPacket(ctx, sourcePort, sourceChannel, destinationPort, destinationChannel, chainType, data)
}

func (k Keeper) CreateOutgoingPacket(
	ctx sdk.Context,
	sourcePort,
	sourceChannel,
	destinationPort,
	destinationChannel,
	chainType string,
	data interface{},
) error {
	if chainType == types.CosmosSdkChainType {
		if data == nil {
			return types.ErrInvalidOutgoingData
		}

		var msgs []sdk.Msg

		switch data := data.(type) {
		case []sdk.Msg:
			msgs = data
		case sdk.Msg:
			msgs = []sdk.Msg{data}
		default:
			return types.ErrInvalidOutgoingData
		}

		interchainAccountTx := types.InterchainAccountTx{Msgs: msgs}

		txBytes, err := k.counterpartyTxCdc.MarshalBinaryBare(interchainAccountTx)
		if err != nil {
			return sdkerrors.Wrap(err, "invalid packet data or codec")
		}

		channelCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(sourcePort, sourceChannel))
		if !ok {
			return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
		}

		// get the next sequence
		sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
		if !found {
			return channel.ErrSequenceSendNotFound
		}

		packetData := types.RunTxPacketData{
			TxBytes: txBytes,
		}

		// TODO: Add timeout height and timestamp
		packet := channel.NewPacket(
			packetData.GetBytes(),
			sequence,
			sourcePort,
			sourceChannel,
			destinationPort,
			destinationChannel,
			math.MaxUint64,
			0,
		)

		return k.channelKeeper.SendPacket(ctx, channelCap, packet)
	} else {
		return sdkerrors.Wrap(types.ErrUnsupportedChainType, chainType)
	}
}

func (k Keeper) DeserializeTx(_ sdk.Context, txBytes []byte) (types.InterchainAccountTx, error) {
	tx := types.InterchainAccountTx{}

	err := k.txCdc.UnmarshalBinaryBare(txBytes, &tx)
	return tx, err
}

func (k Keeper) RunTx(ctx sdk.Context, sourcePort, sourceChannel string, tx types.InterchainAccountTx) error {
	identifier := types.GetIdentifier(sourcePort, sourceChannel)
	err := k.AuthenticateTx(ctx, tx, identifier)
	if err != nil {
		return err
	}

	msgs := tx.Msgs

	// Use cache context.
	// Receive packet msg should succeed regardless of the result of logic.
	// But, if we just return the success even though handler is failed,
	// the leftovers of state transition in handler will remain.
	// However, this can make the unexpected error.
	// To solve this problem, use cache context instead of context,
	// and write the state transition if handler succeeds.
	cacheContext, writeFn := ctx.CacheContext()
	err = nil
	for _, msg := range msgs {
		_, msgErr := k.RunMsg(cacheContext, msg)
		if msgErr != nil {
			err = msgErr
			break
		}
	}

	if err != nil {
		return err
	}

	// Write the state transitions if all handlers succeed.
	writeFn()

	return nil
}

func (k Keeper) AuthenticateTx(ctx sdk.Context, tx types.InterchainAccountTx, identifier string) error {
	msgs := tx.Msgs

	seen := map[string]bool{}
	var signers []sdk.AccAddress
	for _, msg := range msgs {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, addr)
				seen[addr.String()] = true
			}
		}
	}

	store := ctx.KVStore(k.storeKey)

	for _, signer := range signers {
		// Check where the interchain account is made from.
		path := store.Get(signer)
		if !bytes.Equal(path, []byte(identifier)) {
			return sdkerrors.ErrUnauthorized
		}
	}

	return nil
}

func (k Keeper) RunMsg(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
	hander := k.router.Route(ctx, msg.Route())
	if hander == nil {
		return nil, types.ErrInvalidRoute
	}

	return hander(ctx, msg)
}
