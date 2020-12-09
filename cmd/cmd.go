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

// func serverCommand() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "server",
// 		Short: "server",
// 		Long:  "run server",
// 		Run: func(_ *cobra.Command, _ []string) {
// 			server.NewServer()
// 		},
// 	}
// 	return cmd
// }

func Commands(root *cobra.Command, childs ...*cobra.Command) {
	root.AddCommand(
		conversionCommand(),
		uuidCommand(),
		nomadCommand(),
		adduserCommand(),
		// serverCommand(),
	)
	root.Flags().StringVarP(configHelp())
}
