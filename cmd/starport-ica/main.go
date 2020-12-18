package main

import (
	"os"

	"github.com/chainapsis/cosmos-sdk-interchain-account/cmd/starport-ica/cmd"
)

func main() {
	rootCmd := cmd.New()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
