package dispatch

import (
	"context"
	"sync"
	"time"

	"github.com/noaway/v2agent/agent"
	"github.com/noaway/v2agent/agent/core"
	"github.com/noaway/v2agent/config"
	"github.com/noaway/v2agent/internal/worker"
	"github.com/sirupsen/logrus"
)

type User struct {
	core.User
}

type DispatchHandle interface {
	AddUser(*User) error
	DelUser(email string) error
	Close() error
}

type Dispatch struct {
	worker.Guardian
	sync.Once
}

func DispatchStart() DispatchHandle {
	dsp := Dispatch{Guardian: worker.NewGuardian()}
	agent.AgentInit()
	dsp.BreadBoard(dsp.sched()...)
	return &dsp
}

func (dsp *Dispatch) sched() []worker.LineFunc {
	interval := time.Second * 3
	if config.Configure().Agent.SyncInterval > 0 {
		interval = time.Second * time.Duration(config.Configure().Agent.SyncInterval)
	}
	return []worker.LineFunc{
		dsp.do(interval, agent.ContextAgent().SyncUser),
	}
}

func (dsp *Dispatch) do(d time.Duration, fn func() error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		if err := fn(); err != nil {
			return err
		}
		if d == time.Duration(0) {
			return nil
		}
		tick := time.NewTicker(d)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				logrus.Info("sched is closed")
				return nil
			case <-tick.C:
				if err := fn(); err != nil {
					return err
				}
			}
		}
	}
}

func (dsp *Dispatch) AddUser(u *User) error { return agent.ContextAgent().AddUser(&u.User) }

func (dsp *Dispatch) DelUser(email string) error { return agent.ContextAgent().DelUser(email) }

func (dsp *Dispatch) Close() error {
	dsp.Do(func() {
		dsp.Guardian.Close()
	})
	return nil
}
