package cmd

import (
	"github.com/noaway/v2-agent/agent"
	"github.com/noaway/v2-agent/client"
	"github.com/noaway/v2-agent/config"
	"github.com/noaway/v2-agent/internal/svc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type proc struct {
	*config.Configuration
	agent *agent.Agent
}

func (p *proc) Init() error {
	config.NewConfigure(configPath)
	p.Configuration = config.Configure()

	ag, err := agent.NewConfig(
		agent.SetupCluster(
			p.Server.AdvertiseAddr,
			p.Server.BindAddr,
			p.JoinClusterAddrs...,
		),

		agent.SetupUserEventHandler(func(e agent.UserEvent) {

		}),

		agent.SetupDataDir(p.DataDir),
		agent.SetupNodeName(p.Name),
	).NewAgent()

	if err != nil {
		return err
	}

	p.agent = ag
	return nil
}

func (p *proc) Start() error {
	client.DelUser("ggg")
	return nil
}

func (p *proc) Stop() error {
	if p.agent != nil {
		p.agent.Close()
	}
	return nil
}

func webCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "web",
		Short: "Run v2-agnet web",
		Long:  `Run v2-agnet manager web server`,
		Run: func(_ *cobra.Command, _ []string) {
			if err := svc.Run(new(proc), nil); err != nil {
				logrus.Error(err)
			}
		},
	}

	return cmd
}
