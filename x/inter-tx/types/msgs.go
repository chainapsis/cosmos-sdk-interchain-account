package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRegister{}

func NewMsgRegister(sourcePort string, sourceChannel string, sender sdk.AccAddress) *MsgRegister {
	return &MsgRegister{SourcePort: sourcePort, SourceChannel: sourceChannel, Sender: sender}
}

func (MsgRegister) Route() string {
	return RouterKey
}

func (MsgRegister) Type() string {
	return RouterKey
}

func (msg MsgRegister) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	return nil
}

func (msg MsgRegister) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRegister) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

var _ sdk.Msg = &MsgSend{}

func NewMsgSend(sourcePort string, sourceChannel, typ string, amt []sdk.Coin, sender sdk.AccAddress, toAddress sdk.AccAddress) *MsgSend {
	return &MsgSend{SourcePort: sourcePort, SourceChannel: sourceChannel, Typ: typ, Amount: amt, Sender: sender, ToAddress: toAddress}
}

func (MsgSend) Route() string {
	return RouterKey
}

func (MsgSend) Type() string {
	return RouterKey
}

func (msg MsgSend) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if msg.ToAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing recipient address")
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	if !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}

func (msg MsgSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
