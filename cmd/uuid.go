package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/twinj/uuid"
)

func uuidCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uuid",
		Short: "gen uuid",
		Long:  "generate a uuid",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(uuid.NewV4().String())
		},
	}
	return cmd
}
