package cli

import (
	"fmt"

	"github.com/chainapsis/cosmos-sdk-interchain-account/x/inter-tx/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// cmd.AddCommand(GetIBCAccountCmd())

	return cmd
}

func GetIBCAccountCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "ibcaccount [account] [source-port] [source-channel]",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// Verify bech32 address
			acc, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := types.QueryRegisteredParams{Account: acc, SourcePort: args[1], SourceChannel: args[2]}
			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRegistered)

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return fmt.Errorf("failed to marshal params: %w", err)
			}

			res, _, err := clientCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var result []byte
			err = cdc.UnmarshalJSON(res, &result)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(sdk.AccAddress(result).String())
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
