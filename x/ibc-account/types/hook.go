package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type IBCAccountHooks interface {
	// Called when registering IBC account is requested.
	WillAccountCreate(ctx sdk.Context, chainID string, address sdk.AccAddress)
	// Called when IBC account is registered to counterparty chain and acknowledgement packet is delivered.
	DidAccountCreated(ctx sdk.Context, chainID string, address sdk.AccAddress)

	// Called when tx is requested to IBC account.
	WillTxRun(ctx sdk.Context, chainID string, txHash []byte, data interface{})
	// Called when tx is executed and acknowledgement packet is delivered.
	DidTxRun(ctx sdk.Context, chainID string, txHash []byte, data interface{})
}
