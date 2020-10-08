package keeper

import (
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// RegisterIBCAccount performs registering IBC account.
// It will generate the deterministic address by hashing {sourcePort}/{sourceChannel}{salt}.
func (k Keeper) registerIBCAccount(ctx sdk.Context, sourcePort, sourceChannel, destPort, destChannel string, salt []byte) (types.IBCAccountI, error) {
	identifier := types.GetIdentifier(destPort, destChannel)
	address := k.GenerateAddress(identifier, salt)

	account := k.accountKeeper.GetAccount(ctx, address)
	// TODO: Discuss the vulnerabilities when creating a new account only if the old account does not exist
	// Attackers can interrupt creating accounts by sending some assets before the packet is delivered.
	// So it is needed to check that the account is not created from users.
	// Returns an error only if the account was created by other chain.
	// We need to discuss how we can judge this case.
	if account != nil {
		return nil, sdkerrors.Wrap(types.ErrAccountAlreadyExist, account.String())
	}

	ibcAccount := types.NewIBCAccount(
		authtypes.NewBaseAccountWithAddress(address),
		sourcePort, sourceChannel, destPort, destChannel,
	)
	k.accountKeeper.NewAccount(ctx, ibcAccount)
	k.accountKeeper.SetAccount(ctx, ibcAccount)

	return ibcAccount, nil
}

// Determine account's address that will be created.
func (k Keeper) GenerateAddress(identifier string, salt []byte) []byte {
	return tmhash.SumTruncated(append([]byte(identifier), salt...))
}
