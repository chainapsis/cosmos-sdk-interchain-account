package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type IBCAccountHooks interface {
	OnAccountCreated(ctx sdk.Context, chainID string, address sdk.AccAddress)
	OnTxSucceeded(ctx sdk.Context, chainID string, txHash []byte, txBytes []byte)
	OnTxFailed(ctx sdk.Context, chainID string, txHash []byte, txBytes []byte)
}
