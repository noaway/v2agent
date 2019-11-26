package cmd

import (
	"github.com/spf13/cobra"
)

func agentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Run v2-agnet agent",
		Long:  `Run v2-agnet agent manager v2ray CRUD`,
		Run: func(_ *cobra.Command, _ []string) {
				
		},
	}

	return cmd
}
