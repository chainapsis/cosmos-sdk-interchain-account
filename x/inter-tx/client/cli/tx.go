package cli

import (
	"github.com/chainapsis/cosmos-sdk-interchain-account/x/inter-tx/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                types.ModuleName,
		DisableFlagParsing: true,
		RunE:               client.ValidateCmd,
	}

	txCmd.AddCommand(NewRegisterCmd(), NewSendTxCmd())

	return txCmd
}

func NewRegisterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "register",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			sourcePort := viper.GetString(FlagSourcePort)
			sourceChannel := viper.GetString(FlagSourceChannel)

			msg := types.NewMsgRegister(sourcePort, sourceChannel, clientCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsSourcePort)
	cmd.Flags().AddFlagSet(fsSourceChannel)

	_ = cmd.MarkFlagRequired(FlagSourcePort)
	_ = cmd.MarkFlagRequired(FlagSourceChannel)

	return cmd
}

func NewSendTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "send [type] [to_address] [amount]",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			toAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			sourcePort := viper.GetString(FlagSourcePort)
			sourceChannel := viper.GetString(FlagSourceChannel)

			msg := types.NewMsgSend(sourcePort, sourceChannel, args[0], coins, clientCtx.GetFromAddress(), toAddr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsSourcePort)
	cmd.Flags().AddFlagSet(fsSourceChannel)

	_ = cmd.MarkFlagRequired(FlagSourcePort)
	_ = cmd.MarkFlagRequired(FlagSourceChannel)

	return cmd
}
