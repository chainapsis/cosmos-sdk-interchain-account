# Changelog

## PreHistory

##### July, 2020

 * Implement MVP version of interchain account
 * Update cosmos-sdk version to cosmos-sdk#43837b16e768
 * Change the name of interchain account keeper from `InterchainAccountKeeper` to `IBCAccountKeeper`
 * Rename some field from `InterchainAccount` to `IBCAccount`
 * Make `IBCAccountKeeper` handle `onRecvPacket` and `onTimeoutPacket`
 * Make some `try` method like `TryRegisterIBCAccount` and `TryRunTx` able to be used in other keepers
 * Add `hook` concept similar to `staking hook` to help other keepers to communicate with other chains based on `sourcePort` and `sourceChannel`

##### June, 2020

 * Implement PoC version of interchain account
