package types_test

import (
	"encoding/json"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/chainapsis/cosmos-sdk-interchain-account/x/ibc-account/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
)

func TestIBCAccountMarshal(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	baseAcc := authtypes.NewBaseAccountWithAddress(addr)
	ibcAcc := types.NewIBCAccount(baseAcc, "sp", "sc", "dp", "dc")

	// Can't set the seq or pub key on the IBC account
	err := ibcAcc.SetPubKey(pubkey)
	require.Error(t, err)
	err = ibcAcc.SetSequence(1)
	require.Error(t, err)

	bz, err := app.AccountKeeper.MarshalAccount(ibcAcc)
	require.NoError(t, err)

	ibcAcc2, err := app.AccountKeeper.UnmarshalAccount(bz)
	require.NoError(t, err)
	require.Equal(t, ibcAcc.String(), ibcAcc2.String())

	// error on bad bytes
	_, err = app.AccountKeeper.UnmarshalAccount(bz[:len(bz)/2])
	require.Error(t, err)
}

func TestGenesisAccountValidate(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	baseAcc := authtypes.NewBaseAccountWithAddress(addr)

	tests := []struct {
		name   string
		acc    authtypes.GenesisAccount
		expErr bool
	}{
		{
			"valid ibc account",
			types.NewIBCAccount(baseAcc, "sp", "sc", "dp", "dc"),
			false,
		},
		{
			"invalid ibc account that has empty field",
			types.NewIBCAccount(baseAcc, "", "sc", "dp", "dc"),
			true,
		},
		{
			"invalid ibc account that has empty field",
			types.NewIBCAccount(baseAcc, "sp", "", "dp", "dc"),
			true,
		},
		{
			"invalid ibc account that has empty field",
			types.NewIBCAccount(baseAcc, "sp", "sc", "", "dc"),
			true,
		},
		{
			"invalid ibc account that has empty field",
			types.NewIBCAccount(baseAcc, "sp", "sc", "dp", ""),
			true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expErr, tt.acc.Validate() != nil)
		})
	}
}

func TestIBCAccountMarshalYAML(t *testing.T) {
	addr, err := sdk.AccAddressFromHex("0000000000000000000000000000000000000000")
	require.NoError(t, err)
	ba := authtypes.NewBaseAccountWithAddress(addr)

	ibcAcc := types.NewIBCAccount(ba, "sp", "sc", "dp", "dc")

	bs, err := yaml.Marshal(ibcAcc)
	require.NoError(t, err)

	want := "|\n  address: cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a\n  public_key: \"\"\n  account_number: 0\n  sequence: 0\n  source_port: sp\n  source_channel: sc\n  destination_port: dp\n  destination_channel: dc\n"
	require.Equal(t, want, string(bs))
}

func TestIBCAccountJSON(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())
	baseAcc := authtypes.NewBaseAccountWithAddress(addr)

	ibcAcc := types.NewIBCAccount(baseAcc, "sp", "sc", "dp", "dc")

	bz, err := json.Marshal(ibcAcc)
	require.NoError(t, err)

	bz1, err := ibcAcc.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, string(bz), string(bz1))

	var a types.IBCAccount
	require.NoError(t, json.Unmarshal(bz, &a))
	require.Equal(t, ibcAcc.String(), a.String())
}
