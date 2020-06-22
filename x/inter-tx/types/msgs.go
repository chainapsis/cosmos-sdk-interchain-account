package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _ sdk.Msg = MsgRegister{}

func NewMsgRegister(sourcePort string, sourceChannel string, sender sdk.AccAddress) MsgRegister {
	return MsgRegister{SourcePort: sourcePort, SourceChannel: sourceChannel, Sender: sender}
}

func (MsgRegister) Route() string {
	return RouterKey
}

func (MsgRegister) Type() string {
	return RouterKey
}

func (MsgRegister) ValidateBasic() error {
	return nil
}

func (msg MsgRegister) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRegister) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

var _ sdk.Msg = MsgRunTx{}

func NewMsgRunTx(sourcePort string, sourceChannel string, msgBytes []byte, sender sdk.AccAddress) MsgRunTx {
	return MsgRunTx{SourcePort: sourcePort, SourceChannel: sourceChannel, MsgBytes: msgBytes, Sender: sender}
}

func (MsgRunTx) Route() string {
	return RouterKey
}

func (MsgRunTx) Type() string {
	return RouterKey
}

func (MsgRunTx) ValidateBasic() error {
	return nil
}

func (msg MsgRunTx) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRunTx) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
