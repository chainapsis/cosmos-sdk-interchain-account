package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func (MsgTryRegisterIBCAccount) Route() string {
	return ModuleName
}

func (MsgTryRegisterIBCAccount) Type() string {
	return "try-register-ibc-account"
}

func (MsgTryRegisterIBCAccount) ValidateBasic() error {
	return nil
}

func (msg MsgTryRegisterIBCAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgTryRegisterIBCAccount) GetSigners() []sdk.AccAddress {
	// no need to have signer
	return []sdk.AccAddress{msg.Sender}
}

func (MsgTryRunTxMsgSend) Route() string {
	return ModuleName
}

func (MsgTryRunTxMsgSend) Type() string {
	return "try-run-tx-msg-send"
}

func (MsgTryRunTxMsgSend) ValidateBasic() error {
	return nil
}

func (msg MsgTryRunTxMsgSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgTryRunTxMsgSend) GetSigners() []sdk.AccAddress {
	// no need to have signer
	return []sdk.AccAddress{msg.Sender}
}
