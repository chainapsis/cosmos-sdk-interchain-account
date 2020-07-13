package keeper

import (
	"bytes"
	"math"

	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// RegisterIBCAccount performs registering IBC account.
// It will generate the deterministic address by hashing {sourcePort}/{sourceChannel}{salt}.
func (k Keeper) RegisterIBCAccount(
	ctx sdk.Context,
	destPort,
	destChannel,
	salt string,
) error {
	identifier := types.GetIdentifier(destPort, destChannel)
	address := k.GenerateAddress(identifier, salt)
	err := k.CreateAccount(ctx, address, identifier)
	if err != nil {
		return err
	}

	return nil
}

// Create the account if an account with the same address does not exist.
// It will save the address and matched identifier.
func (k Keeper) CreateAccount(ctx sdk.Context, address sdk.AccAddress, identifier string) error {
	account := k.accountKeeper.GetAccount(ctx, address)
	// TODO: Discuss the vulnerabilities when creating a new account only if the old account does not exist
	// Attackers can interrupt creating accounts by sending some assets before the packet is delivered.
	// So it is needed to check that the account is not created from users.
	// Returns an error only if the account was created by other chain.
	// We need to discuss how we can judge this case.
	if account != nil {
		return sdkerrors.Wrap(types.ErrAccountAlreadyExist, account.String())
	}
	// Set account's address if account is nil
	account = k.accountKeeper.NewAccountWithAddress(ctx, address)
	k.accountKeeper.SetAccount(ctx, account)

	store := ctx.KVStore(k.storeKey)
	store = prefix.NewStore(store, types.KeyPrefixRegisteredAccount)
	// Save the identifier for each address to check where the interchain account is made from.
	store.Set(address, []byte(identifier))

	return nil
}

// Determine account's address that will be created.
func (k Keeper) GenerateAddress(identifier string, salt string) []byte {
	return tmhash.SumTruncated([]byte(identifier + salt))
}

// CreateInterchainAccount try to register IBC account to source channel.
// If no source channel exists or doesn't have capability, it will return error.
// Salt is used to generate deterministic address.
func (k Keeper) CreateInterchainAccount(ctx sdk.Context, sourcePort, sourceChannel, salt string) error {
	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrap(channeltypes.ErrChannelNotFound, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return channel.ErrSequenceSendNotFound
	}

	packetData := types.IBCAccountPacketData{
		Type: types.Type_REGISTER,
		Data: []byte(salt),
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

// RequestRunTx try to send messages to source channel.
func (k Keeper) RequestRunTx(ctx sdk.Context, sourcePort, sourceChannel, chainID string, data interface{}) error {
	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrap(channeltypes.ErrChannelNotFound, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	return k.CreateOutgoingPacket(ctx, sourcePort, sourceChannel, destinationPort, destinationChannel, chainID, data)
}

func (k Keeper) CreateOutgoingPacket(
	ctx sdk.Context,
	sourcePort,
	sourceChannel,
	destinationPort,
	destinationChannel,
	chainID string,
	data interface{},
) error {
	if data == nil {
		return types.ErrInvalidOutgoingData
	}

	counterpartyInfo, ok := k.counterpartyInfos[chainID]
	if !ok {
		return types.ErrUnsupportedChain
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

	txBytes, err := counterpartyInfo.SerializeTx(msgs)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid packet data or codec")
	}

	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return channel.ErrSequenceSendNotFound
	}

	packetData := types.IBCAccountPacketData{
		Type: types.Type_RUNTX,
		Data: txBytes,
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
}

func (k Keeper) DeserializeTx(_ sdk.Context, txBytes []byte) ([]sdk.Msg, error) {
	var msgs []sdk.Msg

	err := k.txCdc.UnmarshalBinaryBare(txBytes, &msgs)
	return msgs, err
}

func (k Keeper) RunTx(ctx sdk.Context, destPort, destChannel string, msgs []sdk.Msg) error {
	identifier := types.GetIdentifier(destPort, destChannel)
	err := k.AuthenticateTx(ctx, msgs, identifier)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		err := msg.ValidateBasic()
		if err != nil {
			return err
		}
	}

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

// AuthenticateTx verifies that the messages have the right permission.
// It will check that the message's signers are the IBC account created by the right chain.
func (k Keeper) AuthenticateTx(ctx sdk.Context, msgs []sdk.Msg, identifier string) error {
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
	store = prefix.NewStore(store, types.KeyPrefixRegisteredAccount)

	for _, signer := range signers {
		// Check where the interchain account is made from.
		path := store.Get(signer)
		if !bytes.Equal(path, []byte(identifier)) {
			return sdkerrors.ErrUnauthorized
		}
	}

	return nil
}

// RunMsg executes the message.
// It tries to get the handler from router. And, if router exites, it will perform message.
func (k Keeper) RunMsg(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
	hander := k.router.Route(ctx, msg.Route())
	if hander == nil {
		return nil, types.ErrInvalidRoute
	}

	return hander(ctx, msg)
}

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet) error {
	var data types.IBCAccountPacketData
	// TODO: Remove the usage of global variable "ModuleCdc"
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal interchain account packet data: %s", err.Error())
	}

	switch data.Type {
	case types.Type_REGISTER:
		err := k.RegisterIBCAccount(ctx, packet.DestinationPort, packet.DestinationChannel, string(data.Data))
		if err != nil {
			return err
		}

		return nil
	case types.Type_RUNTX:
		msgs, err := k.DeserializeTx(ctx, data.Data)
		if err != nil {
			return err
		}

		err = k.RunTx(ctx, packet.DestinationPort, packet.DestinationChannel, msgs)
		if err != nil {
			return err
		}

		return nil
	default:
		return types.ErrUnknownPacketData
	}
}
