package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/chainapsis/cosmos-sdk-interchain-account/demo/app"
)

func main() {
	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithTxConfig(encodingConfig.TxConfig).
		WithCodec(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.NewAccountRetriever(encodingConfig.Marshaler)).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome)

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	// TODO: setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc

	rootCmd := &cobra.Command{
		Use:   "democli",
		Short: "Command line interface for interacting with demod",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
			return err
		}

		return server.InterceptConfigsPreRunHandler(cmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCmd(),
		txCmd(),
		flags.LineBreak,
		flags.LineBreak,
		keys.Commands(app.DefaultNodeHome),
		flags.LineBreak,
		version.NewVersionCommand(),
	)

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &client.Context{})
	ctx = context.WithValue(ctx, server.ServerContextKey, server.NewDefaultContext())

	// Add flags and prefix all env exposed with GA
	executor := cli.PrepareMainCmd(rootCmd, "GA", app.DefaultCLIHome)

	err := executor.ExecuteContext(ctx)
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func queryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(),
		flags.LineBreak,
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
		flags.LineBreak,
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd)

	return queryCmd
}

func txCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		bankcmd.NewSendTxCmd(),
		flags.LineBreak,
		authcmd.GetSignCommand(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetMultiSignCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		flags.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd)

	// remove auth and bank commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == authtypes.ModuleName || cmd.Use == banktypes.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}
