package cmd

import "github.com/spf13/cobra"

func New() *cobra.Command {
	c := &cobra.Command{
		Use: "starport-ica",
	}
	c.AddCommand(NewModule())
	return c
}
