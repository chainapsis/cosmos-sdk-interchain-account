package mock

import (
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/testing/mock/keeper"
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/testing/mock/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgTryRegisterIBCAccount:
			return handleMsgTryRegisterIBCAccount(ctx, k, msg)
		case *types.MsgTryRunTxMsgSend:
			return handleMsgTryRunTxMsgSend(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized mock ibc account message type: %T", msg)
		}
	}
}

func handleMsgTryRegisterIBCAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgTryRegisterIBCAccount) (*sdk.Result, error) {
	if err := k.TryRegisterIBCAccount(ctx, msg.SourcePort, msg.SourceChannel, msg.Salt, msg.TimeoutHeight, msg.TimeoutTimestamp); err != nil {
		return nil, err
	}

	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func handleMsgTryRunTxMsgSend(ctx sdk.Context, k keeper.Keeper, msg *types.MsgTryRunTxMsgSend) (*sdk.Result, error) {
	if err := k.TryRunTxMsgSend(ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, msg.FromAddress, msg.ToAddress, msg.Amount); err != nil {
		return nil, err
	}

	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
