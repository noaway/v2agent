package cmd

import (
	"fmt"

	"github.com/noaway/v2agent/config"
	"github.com/noaway/v2agent/core"
	"github.com/spf13/cobra"
)

func testCommand() *cobra.Command {
	var email string
	cmd := &cobra.Command{
		Use:   "test",
		Short: "test",
		Run: func(_ *cobra.Command, _ []string) {
			config.NewConfigure(configPath)
			fmt.Println(core.GetStats(email))
			// fmt.Println(core.Beautify(value))
		},
	}
	cmd.Flags().StringVarP(&configPath, "config", "c", "", "config path")
	cmd.Flags().StringVarP(&email, "email", "", "", "")
	return cmd
}
