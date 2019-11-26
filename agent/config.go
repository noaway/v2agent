package agent

import (
	"io"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/serf"
	"github.com/noaway/v2-agent/internal/utils"
)

const (
	DefaultAdvertisePort = 5421
	DefaultBindPort      = 5422
)

type UserEvent struct{ serf.UserEvent }

type UserEventHandler func(UserEvent)

func SetupCluster(advertiseAddr, bindAddr string, clusterAddrs ...string) ConfOptHandle {
	return func(c *Config) {
		advertiseHost, advertisePort, err := utils.ParseIPAndPort(advertiseAddr)
		if err != nil {
			advertiseHost = "0.0.0.0"
			advertisePort = DefaultAdvertisePort
		}
		c.serfConfig.MemberlistConfig.AdvertiseAddr = advertiseHost
		c.serfConfig.MemberlistConfig.AdvertisePort = advertisePort

		bindHost, bindPort, err := utils.ParseIPAndPort(bindAddr)
		if err != nil {
			bindHost = "0.0.0.0"
			bindPort = DefaultBindPort
		}
		c.serfConfig.MemberlistConfig.BindAddr = bindHost
		c.serfConfig.MemberlistConfig.BindPort = bindPort

		c.ClusterAddrs = clusterAddrs
	}
}

func SetupUserEventHandler(handle UserEventHandler) ConfOptHandle {
	return func(c *Config) { c.UserEventHandler = handle }
}

func SetupDataDir(dataDir string) ConfOptHandle {
	return func(c *Config) { c.DataDir = dataDir }
}

func SetupNodeName(nodeName string) ConfOptHandle {
	return func(c *Config) { c.serfConfig.NodeName = nodeName }
}

type ConfOptHandle func(*Config)

func NewConfig(opts ...ConfOptHandle) *Config {
	c := &Config{}
	serfConf := serf.DefaultConfig()
	{
		serfConf.QueueDepthWarning = 1000000
		serfConf.MinQueueDepth = 4096
		serfConf.LeavePropagateDelay = 3 * time.Second
		serfConf.ReconnectTimeout = 3 * 24 * time.Hour
		serfConf.MemberlistConfig = memberlist.DefaultWANConfig()
		serfConf.MemberlistConfig.AdvertisePort = DefaultAdvertisePort
		serfConf.MemberlistConfig.BindPort = DefaultBindPort
		serfConf.MemberlistConfig.DeadNodeReclaimTime = 30 * time.Second

		serfConf.CoalescePeriod = time.Second * 3
		serfConf.QuiescentPeriod = time.Second * 2
	}
	c.serfConfig = serfConf

	for _, opt := range opts {
		opt(c)
	}
	return c
}

type Config struct {
	serfConfig *serf.Config

	BindAddr         string
	LogOutput        io.Writer
	DataDir          string
	ClusterAddrs     []string
	UserEventHandler UserEventHandler
}

func (conf *Config) NewAgent() (*Agent, error) { return NewAgent(conf) }
