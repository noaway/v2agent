package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/serf/serf"
	"github.com/noaway/v2agent/agent/core"
	"github.com/noaway/v2agent/config"
	"github.com/noaway/v2agent/internal/utils"
	"github.com/sirupsen/logrus"
)

const (
	serfEventBacklogWarning = 200
)

const (
	USERS_BUCKET = "users"
)

const (
	ADD_USER_EVENT    = "add_user_event"
	DELETE_USER_EVENT = "delete_user_event"
)

var (
	onceAgent *AgentIns
	once      = sync.Once{}
)

type Agent interface {
	SyncUser() error
	AddUser(*core.User) error
	DelUser(string) error
}

func AgentInit(fns ...UserEventHandler) {
	once.Do(func() {
		conf := config.Configure().Agent
		config := NewConfig(
			SetupCluster(
				conf.AdvertisePort,
				conf.BindAddr,
				conf.JoinClusterAddrs...,
			),
			SetupDataDir(conf.DataDir),
			SetupNodeName(conf.Name),
			SetupRegion(conf.Region),
		)

		if len(fns) > 0 {
			config.UserEventHandler = fns[0]
		}

		if config.DataDir == "" {
			panic("DataDir is empty")
		}

		if config.Region == "" {
			config.Region = "default"
		}

		if _, err := os.Stat(config.DataDir); err != nil {
			err := os.MkdirAll(config.DataDir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		db, err := bolt.Open(config.DataDir+"/agent.db", 0600, &bolt.Options{Timeout: 5 * time.Second})
		if err != nil {
			panic(fmt.Sprintf("bolt init err %v", err))
		}

		agent := &AgentIns{
			eventCh:    make(chan serf.Event, 256),
			shutdownCh: make(chan struct{}),
			config:     config,
			DB:         db,
		}

		agent.serf, err = agent.createSerf()
		if err != nil {
			panic(err)
		}

		agent.Join(agent.config.ClusterAddrs...)
		go agent.eventHandler()
		onceAgent = agent
	})
}

func ContextAgent() Agent { return onceAgent }

type AgentIns struct {
	*bolt.DB

	config     *Config
	serf       *serf.Serf
	eventCh    chan serf.Event
	shutdownCh chan struct{}
}

func (agent *AgentIns) createSerf() (*serf.Serf, error) {
	conf := agent.config.serfConfig
	conf.Init()

	conf.Tags["name"] = conf.NodeName
	logger := log.New(logrus.StandardLogger().Out, "", log.LstdFlags)
	conf.EventCh = agent.eventCh
	conf.Logger = logger
	conf.LogOutput = logrus.StandardLogger().Out

	conf.SnapshotPath = filepath.Join(agent.config.DataDir + "/v2-agent.log")
	if err := utils.EnsurePath(conf.SnapshotPath, false); err != nil {
		return nil, err
	}

	return serf.Create(conf)
}

func (agent *AgentIns) Join(addrs ...string) (int, error) {
	return agent.serf.Join(addrs, true)
}

func (agent *AgentIns) SyncUser() error {
	if agent == nil {
		return nil
	}

	return agent.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(USERS_BUCKET))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			u := core.User{}
			if err := json.Unmarshal(v, &u); err != nil {
				return err
			}
			if _, err := core.AddUser(&u); err != nil {
				if strings.Contains(err.Error(), "already exists") {
					logrus.Debugf("local user data synchronization to v2ray [email='%v %v']", u.Email, "already exists")
					continue
				}
				return err
			}
			logrus.Debugf("sync user func find new user [user='%v']", u)
		}
		return nil
	})
}

func (agent *AgentIns) AddUser(user *core.User) error {
	data, err := user.Encode()
	if err != nil {
		return err
	}
	return agent.userEvent(ADD_USER_EVENT, data, true)
}

func (agent *AgentIns) DelUser(email string) error {
	return agent.userEvent(DELETE_USER_EVENT, []byte(email), true)
}

func (agent *AgentIns) UserEventHandler(ue UserEvent) {
	entry := logrus.WithFields(logrus.Fields{
		"event type": "event messages received by other servers",
	})
	switch ue.Name {
	case ADD_USER_EVENT:
		u := core.User{}
		if err := u.Decode(ue.Payload); err != nil {
			entry.Warningf("ADD_USER_EVENT decoding [err='%v']", err)
			return
		}
		if agent.InRegion(u.Regions...) {
			agent.saveUser(&u)
		}
	case DELETE_USER_EVENT:
		if err := agent.removeUser(string(ue.Payload)); err != nil {
			entry.Warningf("DELETE_USER_EVENT remove user [err='%v', playload='%v']", err, string(ue.Payload))
		}
	}
}

func (agent *AgentIns) saveUser(u *core.User) error {
	logrus.Infof(
		"save user [uuid='%v', email='%v', region='%v', alterId='%v']",
		u.UUID, u.Email, u.Regions, u.AlterId,
	)
	return agent.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(USERS_BUCKET))
		if err != nil {
			logrus.Warningf("save the user to CreateBucketIfNotExists [err='%v']", err)
			return err
		}
		err = b.Put([]byte(u.Email), u.Data())
		if err != nil {
			logrus.Warningf("save the user to the local db [err='%v']", err)
		}
		return err
	})
}

func (agent *AgentIns) removeUser(email string) error {
	return agent.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(USERS_BUCKET))
		if err != nil {
			return err
		}

		if err := b.Delete([]byte(email)); err != nil {
			return err
		}

		return core.DelUser(email)
	})
}

func (agent *AgentIns) Members() []serf.Member {
	return agent.serf.Members()
}

func (agent *AgentIns) userEvent(name string, payload []byte, coalesce bool) error {
	return agent.serf.UserEvent(name, payload, coalesce)
}

func (agent *AgentIns) eventHandler() {
	var numQueuedEvents int
	for {
		numQueuedEvents = len(agent.eventCh)
		if numQueuedEvents > serfEventBacklogWarning {
			logrus.Warnf("v2agent: number of queued serf events above warning threshold: %d/%d", numQueuedEvents, serfEventBacklogWarning)
		}

		select {
		case e := <-agent.eventCh:
			switch e.EventType() {
			case serf.EventMemberJoin:
				// _ = serf.Member
			case serf.EventMemberLeave, serf.EventMemberFailed, serf.EventMemberReap:
			case serf.EventUser:
				ue := UserEvent{UserEvent: e.(serf.UserEvent)}
				if agent.config.UserEventHandler != nil {
					agent.config.UserEventHandler(ue)
				}
				agent.UserEventHandler(ue)
			case serf.EventMemberUpdate: // Ignore
			case serf.EventQuery: // Ignore
			default:
				logrus.Warnf("consul: unhandled LAN Serf Event: %#v", e)
			}
		case <-agent.shutdownCh:
			return
		}
	}
}

func (agent *AgentIns) InRegion(regions ...string) bool {
	for _, region := range regions {
		if agent.config.Region == region {
			return true
		}
	}
	return false
}

func Close() error {
	if onceAgent != nil {
		if onceAgent.DB != nil {
			onceAgent.DB.Close()
		}
		if onceAgent.serf != nil {
			onceAgent.serf.Leave()
		}
		close(onceAgent.shutdownCh)
	}
	return nil
}
