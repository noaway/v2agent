package agent

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/hashicorp/serf/serf"
	"github.com/noaway/v2-agent/internal/utils"
)

const (
	serfEventBacklogWarning = 200
)

func NewAgent(config *Config) (*Agent, error) {
	if config.DataDir == "" {
		return nil, errors.New("DataDir is empty")
	}

	agent := &Agent{
		eventCh:    make(chan serf.Event, 256),
		shutdownCh: make(chan struct{}),
		config:     config,
	}

	var err error
	agent.serf, err = agent.createSerf("./v2-agent.noaway")
	if err != nil {
		return nil, err
	}

	go agent.eventHandler()
	return agent, nil
}

type Agent struct {
	config     *Config
	logger     *log.Logger
	serf       *serf.Serf
	eventCh    chan serf.Event
	shutdownCh chan struct{}
}

func (agent *Agent) createSerf(path string) (*serf.Serf, error) {
	conf := agent.config.serfConfig
	conf.Init()

	conf.Tags["name"] = conf.NodeName
	if agent.logger == nil {
		conf.MemberlistConfig.LogOutput = agent.config.LogOutput
		conf.LogOutput = agent.config.LogOutput
	}

	conf.Logger = agent.logger
	conf.EventCh = agent.eventCh
	conf.MemberlistConfig.Logger = agent.logger

	conf.SnapshotPath = filepath.Join(agent.config.DataDir, path)
	if err := utils.EnsurePath(conf.SnapshotPath, false); err != nil {
		return nil, err
	}

	return serf.Create(conf)
}

func (agent *Agent) Join(addrs ...string) (int, error) {
	return agent.serf.Join(addrs, true)
}

func (agent *Agent) Members() []serf.Member {
	return agent.serf.Members()
}

func (agent *Agent) UserEvent(name string, payload []byte, coalesce bool) error {
	return agent.serf.UserEvent(name, payload, coalesce)
}

func (agent *Agent) eventHandler() {
	var numQueuedEvents int
	for {
		numQueuedEvents = len(agent.eventCh)
		if numQueuedEvents > serfEventBacklogWarning {
			agent.logger.Printf("[WARN] v2agent: number of queued serf events above warning threshold: %d/%d", numQueuedEvents, serfEventBacklogWarning)
		}

		select {
		case e := <-agent.eventCh:
			switch e.EventType() {
			case serf.EventMemberJoin:
			case serf.EventMemberLeave, serf.EventMemberFailed, serf.EventMemberReap:
			case serf.EventUser:
				if agent.config.UserEventHandler != nil {
					agent.config.UserEventHandler(UserEvent{UserEvent: e.(serf.UserEvent)})
				}
			case serf.EventMemberUpdate: // Ignore
			case serf.EventQuery: // Ignore
			default:
				agent.logger.Printf("[WARN] consul: unhandled LAN Serf Event: %#v", e)
			}
		case <-agent.shutdownCh:
			return
		}
	}
}

func (a *Agent) Close() {
	if a.serf != nil {
		a.serf.Leave()
	}
	close(a.shutdownCh)
}
