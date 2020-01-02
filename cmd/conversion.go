package cmd

import (
	"fmt"
	"github.com/noaway/v2agent/internal/gensub"

	"github.com/noaway/v2agent/config"
	"github.com/spf13/cobra"
)

func getKitsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kits",
		Short: "v2agnet conversion support of kit",
		Long:  `v2agnet conversion support of kit`,
		Run: func(_ *cobra.Command, _ []string) {
			for k := range gensub.KitMap {
				fmt.Println(k)
			}
		},
	}

	return cmd
}

func conversionCommand() *cobra.Command {
	var kitKey string
	cmd := &cobra.Command{
		Use:   "conversion",
		Short: "v2agnet conversion config",
		Long: `unified v2ray configuration file 
will be transformed into different client configuration, 
and finally upload the server to realize the subscription function`,
		Run: func(_ *cobra.Command, _ []string) {
			config.NewConfigure(configPath)
			defer config.Close()

			kit, ok := gensub.KitMap[kitKey]
			if !ok {
				fmt.Println("not found kit")
				return
			}
			fmt.Println(kit.Content(config.Configure().V2CliConfig))
		},
	}
	cmd.Flags().StringVarP(configHelp())
	cmd.Flags().StringVarP(&kitKey, "kit", "", "", "kit")
	cmd.AddCommand(getKitsCommand())
	return cmd
}
