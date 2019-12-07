package cmd

import (
	"github.com/noaway/v2agent/config"
	"github.com/noaway/v2agent/dispatch"
	"github.com/noaway/v2agent/internal/svc"
	"github.com/spf13/cobra"
)

func agentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Run v2-agnet agent",
		Long:  `Run v2-agnet agent manager v2ray CRUD`,
		Run: func(_ *cobra.Command, _ []string) {
			svc.Proc(func(p *svc.Pair) error {
				config.NewConfigure(configPath)
				dsp := dispatch.DispatchStart()
				p.Set("dsp", dsp)
				return nil
			}, func(p *svc.Pair) error {
				dsp := p.Get("dsp").(*dispatch.Dispatch)
				dsp.Close()
				return nil
			})
		},
	}

	cmd.Flags().StringVarP(configHelp())
	return cmd
}
