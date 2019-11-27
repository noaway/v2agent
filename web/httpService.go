package web

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// HTTPService struct
type HTTPService struct {
	sync.RWMutex
	srv *http.Server
}

// Service func
func (hs *HTTPService) Service(listener net.Listener, handler http.Handler) error {
	srv := &http.Server{
		Handler:        handler,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	hs.Lock()
	hs.srv = srv
	hs.Unlock()
	return srv.Serve(listener)
}

// Shutdown func
func (hs *HTTPService) Shutdown() {
	logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := hs.srv.Shutdown(ctx); err != nil {
		logrus.Error("Server Shutdown:", err)
	}
	logrus.Info("Server exiting")
}
