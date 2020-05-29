package cmd

import (
	"github.com/spf13/cobra"
)

func adduserCommand() *cobra.Command {
	cmd:=&cobra.Command{
		Use:   "adduser",
		Short: "adduser",
		Long:  "add a user",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("add user")
		}
	}
	return cmd
}
