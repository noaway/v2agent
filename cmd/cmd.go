package cmd

import (
	"github.com/spf13/cobra"
)

var configPath string

func Commands(root *cobra.Command, childs ...*cobra.Command) {
	root.AddCommand(
		webCommand(),
		agentCommand(),
	)
	root.Flags().StringVarP(&configPath, "config", "c", "", "config path")
}
