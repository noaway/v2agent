package web

import (
	"net"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/noaway/v2agent/dispatch"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	userkey = "user"
)

// NewAPI func
func NewWEB(listener net.Listener) *WEB { return &WEB{HTTPListener: listener} }

// API struct
type WEB struct {
	HTTPService
	HTTPListener net.Listener

	serverID string
}

// Handler func
func (web *WEB) Handler() (http.Handler, error) {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.HandleMethodNotAllowed = true

	r.Any("ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	pprof.Register(r)
	// metrics
	r.GET("metrics", func(c *gin.Context) { promhttp.Handler().ServeHTTP(c.Writer, c.Request) })

	dsp := dispatch.DispatchStart()
	api := r.Group("/v1/api")
	{
		api.POST("user", addUser(dsp))
		api.DELETE("user", delUser(dsp))
	}

	return r, nil
}

// Main func
func (web *WEB) Main() error {
	r, err := web.Handler()
	if err != nil {
		return err
	}
	if err := web.Service(web.HTTPListener, r); err != nil {
		web.Close()
	}
	return err
}

// Close func
func (web *WEB) Close() {
	web.RLock()
	defer web.RUnlock()

	if web.srv != nil {
		web.Shutdown()
	}
}
