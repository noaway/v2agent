package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configPath string
)

func configHelp() (p *string, name, shorthand string, value string, usage string) {
	p = &configPath
	name = "config"
	shorthand = "c"
	usage = "global config"

	return
}

func Commands(root *cobra.Command, childs ...*cobra.Command) {
	root.AddCommand(
		webCommand(),
		agentCommand(),
		conversionCommand(),
		uuidCommand(),
	)
	root.Flags().StringVarP(configHelp())
}
