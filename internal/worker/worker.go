package worker

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/noaway/v2agent/internal/utils"
)

// Mould struct
type Mould func(...interface{})

type worker struct {
	Guardian

	name      string
	workerCnt int
	m         Mould
	channel   chan []interface{}
}

func (w *worker) run(num int) {
	w.BreadBoard(func(ctx context.Context) error {
		for args := range w.channel {
			w.m(args...)
		}
		return nil
	})
}

// Workers struct
type Workers struct {
	sync.RWMutex

	channels map[string]*worker
}

// NewWorkers func
func NewWorkers() *Workers {
	return &Workers{channels: make(map[string]*worker)}
}

// Run func
func (ws *Workers) run(w *worker) {
	for i := 0; i < w.workerCnt; i++ {
		go w.run(i)
	}
}

// HandleFunc func
func (ws *Workers) HandleFunc(name string, m Mould, workerCnt int) {
	ws.Lock()
	defer ws.Unlock()
	if _, ok := ws.channels[name]; !ok {
		w := &worker{
			name:      name,
			workerCnt: workerCnt,
			m:         m,
			channel:   make(chan []interface{}, workerCnt*2),
			Guardian:  newGuardian(),
		}
		ws.run(w)
		ws.channels[name] = w
	}
}

// Remove func
func (ws *Workers) Remove(name string) {
	ws.Lock()
	defer ws.Unlock()
	if worker, ok := ws.channels[name]; ok {
		close(worker.channel)
		delete(ws.channels, name)
	}
}

// Transmit func
func (ws *Workers) Transmit(name string, args ...interface{}) {
	ws.RLock()
	defer ws.RUnlock()
	if worker, ok := ws.channels[name]; ok {
		worker.channel <- args
	}
}

// Close func
func (ws *Workers) Close() {
	ws.Lock()
	defer ws.Unlock()
	for name, worker := range ws.channels {
		close(worker.channel)
		delete(ws.channels, name)
	}
}

func Go(fn interface{}, args ...interface{}) {
	nf := utils.NewFunction(fn)
	go func() {
		defer utils.DeferError(func(stack string, err interface{}) {
			logrus.Errorf("Go defer err: %v stack info: %v", err, stack)
		})
		nf.Invoke(args...)
	}()
}
