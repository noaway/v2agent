package svc

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/noaway/v2agent/internal/utils"
	"github.com/sirupsen/logrus"
)

var signalNotify = signal.Notify

// Service interface
type Service interface {
	Init() error
	Start() error
	Stop() error
}

// Run run a service
func Run(service Service, signalFunc func() error, sig ...os.Signal) error {
	defer utils.DeferError(func(stack string, err interface{}) {
		logrus.Errorf("Run.trace err: %v stack info: %v", err, stack)
	})

	if err := service.Init(); err != nil {
		return err
	}

	if err := service.Start(); err != nil {
		return err
	}

	if signalFunc == nil {
		if len(sig) == 0 {
			sig = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
		}

		signalChan := make(chan os.Signal, 1)
		signalNotify(signalChan, sig...)
		<-signalChan
		close(signalChan)
	} else {
		if err := signalFunc(); err != nil {
			return err
		}
	}

	return service.Stop()
}

type Pair struct {
	m map[string]interface{}
}

func (p *Pair) Get(k string) interface{}    { return p.m[k] }
func (p *Pair) Set(k string, v interface{}) { p.m[k] = v }

func Proc(begin func(*Pair) error, end func(*Pair) error) error {
	p := &Pair{make(map[string]interface{})}
	if err := begin(p); err != nil {
		return err
	}
	sig := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	signalChan := make(chan os.Signal, 1)
	signalNotify(signalChan, sig...)
	<-signalChan
	close(signalChan)
	return end(p)
}

// BaseWrapper struct
type BaseWrapper struct {
	sync.WaitGroup
	sync.Once
}

// Go func
func (w *BaseWrapper) Go(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}
