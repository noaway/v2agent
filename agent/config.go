package agent

import (
	"io"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/hashicorp/serf/serf"
	"github.com/noaway/v2agent/internal/utils"
)

const (
	DefaultAdvertiseBindPort = 5421
	DefaultBindPort          = 5422
)

type UserEventHandler func(UserEvent)

type UserEvent struct{ serf.UserEvent }

func SetupCluster(advertisePort int, bindAddr string, clusterAddrs ...string) ConfOptHandle {
	return func(c *Config) {
		if advertisePort == 0 {
			advertisePort = DefaultAdvertiseBindPort
		}
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

func SetupRegion(region string) ConfOptHandle {
	return func(c *Config) { c.Region = region }
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
		serfConf.MemberlistConfig = memberlist.DefaultLANConfig()
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
	serfConfig       *serf.Config
	BindAddr         string
	LogOutput        io.Writer
	DataDir          string
	ClusterAddrs     []string
	Region           string
	UserEventHandler UserEventHandler
}
