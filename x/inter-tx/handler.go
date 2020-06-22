package inter_tx

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case MsgRegister:
			return handleMsgRegister(ctx, msg, k)
		case MsgRunTx:
			return handleMsgRunTx(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgRegister(ctx sdk.Context, msg MsgRegister, k Keeper) (*sdk.Result, error) {
	err := k.RegisterInterchainAccount(ctx, msg.Sender, msg.SourcePort, msg.SourceChannel)

	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgRunTx(ctx sdk.Context, msg MsgRunTx, k Keeper) (*sdk.Result, error) {
	err := k.RunMsg(ctx, msg.MsgBytes, msg.Sender, msg.SourcePort, msg.SourceChannel)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
