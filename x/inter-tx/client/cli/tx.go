package cli

import (
	"bufio"
	"fmt"

	"github.com/chainapsis/cosmos-sdk-interchain-account/x/inter-tx/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	interTxTxCmd := &cobra.Command{
		Use:                types.ModuleName,
		DisableFlagParsing: true,
		RunE:               client.ValidateCmd,
	}

	interTxTxCmd.AddCommand(flags.PostCommands(GetRegisterCmd(cdc), GetSendTxCmd(cdc))...)

	return interTxTxCmd
}

func GetRegisterCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "register",
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			txBldr, msg, err := BuildRegisterMsg(cliCtx, txBldr)
			if err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsSourcePort)
	cmd.Flags().AddFlagSet(fsSourceChannel)

	_ = cmd.MarkFlagRequired(FlagSourcePort)
	_ = cmd.MarkFlagRequired(FlagSourceChannel)

	return cmd
}

func GetSendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "send [to_address] [amount]",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			txBldr, msg, err := BuildSendTxMsg(cliCtx, txBldr, args)
			if err != nil {
				return err
			}
			return authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsSourcePort)
	cmd.Flags().AddFlagSet(fsSourceChannel)

	_ = cmd.MarkFlagRequired(FlagSourcePort)
	_ = cmd.MarkFlagRequired(FlagSourceChannel)

	return cmd
}

func BuildRegisterMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder) (auth.TxBuilder, sdk.Msg, error) {
	sender := cliCtx.GetFromAddress()

	sourcePort := viper.GetString(FlagSourcePort)
	sourceChannel := viper.GetString(FlagSourceChannel)

	msg := types.NewMsgRegister(sourcePort, sourceChannel, sender)

	return txBldr, msg, nil
}

func BuildSendTxMsg(cliCtx context.CLIContext, txBldr auth.TxBuilder, args []string) (auth.TxBuilder, sdk.Msg, error) {
	sourcePort := viper.GetString(FlagSourcePort)
	sourceChannel := viper.GetString(FlagSourceChannel)

	sender := cliCtx.GetFromAddress()

	to, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return txBldr, nil, err
	}

	// parse coins trying to be sent
	coins, err := sdk.ParseCoins(args[1])
	if err != nil {
		return txBldr, nil, err
	}

	// Get ibc account
	params := types.QueryRegisteredParams{Account: sender, SourcePort: sourcePort, SourceChannel: sourceChannel}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRegistered)

	bz, err := cliCtx.Codec.MarshalJSON(params)
	if err != nil {
		return txBldr, nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	res, _, err := cliCtx.QueryWithData(route, bz)
	if err != nil {
		return txBldr, nil, err
	}

	var result []byte
	err = cliCtx.Codec.UnmarshalJSON(res, &result)
	if err != nil {
		return txBldr, nil, err
	}
	ibcAccount := sdk.AccAddress(result)

	msg := banktypes.NewMsgSend(ibcAccount, to, coins)
	bz, err = cliCtx.Codec.MarshalBinaryBare(msg)
	if err != nil {
		return txBldr, nil, err
	}

	return txBldr, types.NewMsgRunTx(sourcePort, sourceChannel, bz, sender), nil
}
