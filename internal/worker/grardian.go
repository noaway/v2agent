package worker

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/noaway/v2agent/internal/utils"
)

type LineFunc func(context.Context) error

// Guardian struct
type Guardian struct {
	sync.WaitGroup

	gctx    context.Context
	gcannel context.CancelFunc
	timeout time.Duration
}

func NewGuardian() Guardian {
	return newGuardian()
}

func newGuardian() Guardian {
	ctx, cancel := context.WithCancel(context.Background())
	return Guardian{gctx: ctx, gcannel: cancel}
}

// BreadBoard func
func (g *Guardian) BreadBoard(lines ...LineFunc) {
	for _, line := range lines {
		if line == nil {
			continue
		}
		g.Add(1)
		go g.Run(line)
	}
}

// Trace func
func (g *Guardian) Trace(v ...interface{}) {
	logrus.Info(v)
}

// Run func
func (g *Guardian) Run(line LineFunc) {
	for {
		select {
		case <-g.gctx.Done():
			return
		default:
			err := func() error {
				defer utils.DeferError(func(stack string, err interface{}) {
					logrus.Errorf("Guardian.Run.trace err: %v stack info: %v", err, stack)
				})

				if err := line(g.gctx); err != nil {
					logrus.Error("Guardian.return.line err: ", err)
					return err
				}
				return nil
			}()
			if err != nil {
				logrus.Info("2 seconds later retry err ", err)
			}
			time.Sleep(time.Second * 2)
		}
	}
}

func (g *Guardian) Close() {
	g.gcannel()
}
