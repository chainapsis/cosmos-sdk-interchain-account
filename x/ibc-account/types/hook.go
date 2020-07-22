package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type IBCAccountHooks interface {
	OnAccountCreated(ctx sdk.Context, sourcePort, sourceChannel string, address sdk.AccAddress)
	OnTxSucceeded(ctx sdk.Context, sourcePort, sourceChannel string, txHash []byte, txBytes []byte)
	OnTxFailed(ctx sdk.Context, sourcePort, sourceChannel string, txHash []byte, txBytes []byte)
}
