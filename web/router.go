package web

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
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
	var files = []string{}
	filepath.Walk("template/default", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			files = append(files, path)
		}
		return nil
	})
	logrus.Info(files)
	r.LoadHTMLFiles(files...)
	r.Static("/public", "template/default")

	r.Use(sessions.Sessions("mysession", cookie.NewStore([]byte("noaway_loginuser"))))

	gin.SetMode(gin.ReleaseMode)
	r.HandleMethodNotAllowed = true

	r.Any("ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	pprof.Register(r)
	// metrics
	r.GET("metrics", func(c *gin.Context) { promhttp.Handler().ServeHTTP(c.Writer, c.Request) })

	v1 := r.Group("v1", authRequired)
	{
		v1.GET("getips", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", "Hello World noaway")
		})

	}

	auth := r.Group("auth")
	{
		auth.GET("/login", getLogin)
		auth.POST("/login", postLogin)
		auth.GET("/register", getRegister)
		auth.POST("/register", postRegister)
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
