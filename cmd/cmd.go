package cmd

import (
	"github.com/spf13/cobra"
)

var configPath string

func Commands(root *cobra.Command, childs ...*cobra.Command) {
	root.Flags().StringVarP(&configPath, "config", "c", "", "config path")

	root.AddCommand(
		webCommand(),
		agentCommand(),
	)
}
