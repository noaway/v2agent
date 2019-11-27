package web

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
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

	store := cookie.NewStore([]byte("loginuser"))
    r.Use(sessions.Sessions("mysession", store))

	gin.SetMode(gin.ReleaseMode)
	r.HandleMethodNotAllowed = true

	r.Any("ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	pprof.Register(r)
	// metrics
	r.GET("metrics", func(c *gin.Context) { promhttp.Handler().ServeHTTP(c.Writer, c.Request) })

	v1 := r.Group("v1")
	{
		v1.GET("getips", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", "Hello World noaway")
		})

		v1.GET("login", func(c *gin.Context) {
			  
			// session := sessions.Default(c)
			// session.Set("loginuser", username)
			// session.Save()
			
			c.HTML(http.StatusOK, "login.html", "Hello World noaway")
		})
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
