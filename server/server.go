package server

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"net"
// 	"net/http"
// 	"sync"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/noaway/v2agent/config"
// 	"github.com/noaway/v2agent/internal/svc"
// 	"github.com/noaway/v2agent/models"
// 	"github.com/sirupsen/logrus"
// )

// // HTTPService struct
// type HTTPService struct {
// 	sync.RWMutex
// 	srv *http.Server
// }

// // Service func
// func (hs *HTTPService) Service(listener net.Listener, handler http.Handler) error {
// 	srv := &http.Server{
// 		Handler:        handler,
// 		ReadTimeout:    30 * time.Second,
// 		WriteTimeout:   30 * time.Second,
// 		MaxHeaderBytes: 1 << 20,
// 	}
// 	hs.Lock()
// 	hs.srv = srv
// 	hs.Unlock()
// 	return srv.Serve(listener)
// }

// // Shutdown func
// func (hs *HTTPService) Shutdown() {
// 	logrus.Info("Shutdown Server ...")

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	if err := hs.srv.Shutdown(ctx); err != nil {
// 		logrus.Error("Server Shutdown:", err)
// 	}
// 	logrus.Info("Server exiting")
// }

// func NewServer(configPath string) error {
// 	flag.Parse()
// 	if err := config.NewConfigure(configPath); err != nil {
// 		return err
// 	}

// 	if err := models.InitPostgre(); err != nil {
// 		return fmt.Errorf("InitPostgre err %v", err)
// 	}

// 	httpListener, err := net.Listen("tcp", config.Configure.Server.Addr)
// 	if err != nil {
// 		return err
// 	}
// 	server := HTTPService{}
// 	r := gin.New()
// 	gin.SetMode(gin.ReleaseMode)
// 	server.Service(httpListener, r)
// 	svc.WaitSignal()
// 	server.Shutdown()
// }
