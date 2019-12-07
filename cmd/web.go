package cmd

import (
	"fmt"
	"net"

	"github.com/noaway/v2agent/agent"
	// "github.com/noaway/v2agent/client"
	"github.com/noaway/v2agent/config"
	"github.com/noaway/v2agent/internal/svc"
	"github.com/noaway/v2agent/web"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type proc struct {
	svc.BaseWrapper
	*config.Configuration

	agent *agent.Agent
	web   *web.WEB
}

func (p *proc) Init() error {
	config.NewConfigure(configPath)
	p.Configuration = config.Configure()

	// ag, err := agent.NewConfig(
	// 	agent.SetupCluster(
	// 		p.Server.AdvertiseAddr,
	// 		p.Server.BindAddr,
	// 		p.JoinClusterAddrs...,
	// 	),

	// 	agent.SetupUserEventHandler(func(e agent.UserEvent) {

	// 	}),

	// 	agent.SetupDataDir(p.DataDir),
	// 	agent.SetupNodeName(p.Name),
	// ).NewAgent()

	// if err != nil {
	// 	return err
	// }

	// p.agent = ag
	return nil
}

func (p *proc) Start() error {
	// client.DelUser("ggg")
	if p.Server.HttpAddr == "" {
		return fmt.Errorf("%v", "[web] error listen addr is empty")
	}
	logrus.Info("listen: ", p.Server.HttpAddr)
	httpListener, err := net.Listen("tcp", p.Server.HttpAddr)
	if err != nil {
		return err
	}
	p.web = web.NewWEB(httpListener)
	p.Go(func() { p.web.Main() })

	return nil
}

func (p *proc) Stop() error {
	if p.web != nil {
		p.web.Close()
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
	cmd.Flags().StringVarP(&configPath, "config", "c", "", "config path")
	return cmd
}
