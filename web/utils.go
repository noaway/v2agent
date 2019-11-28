package web

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	DefaultKey = "github/noaway/v2agent/web/viewUser"
)

type Page string

func (p Page) String() string { return string(p) }

// page
const (
	LOGIN Page = "login.html"
)

func Ctx(c *gin.Context) *Context {
	ctx := c.MustGet(DefaultKey).(*Context)
	if ctx == nil {
		c.Abort()
		return nil
	}
	ctx.Session = sessions.Default(c)
	if ctx.User == nil {
		ctx.User = new(User)
	}
	return ctx
}

func (c *Context) View(page Page, vs ...gin.H) {
	obj := gin.H{"User": c.User}
	for _, v := range vs {
		for k, v := range v {
			if k == "User" {
				continue
			}
			obj[k] = v
		}
	}
	c.HTML(http.StatusOK, page.String(), obj)
}

type Context struct {
	*gin.Context

	Session sessions.Session
	User    *User
}

type User struct {
	UserName string
}

func authRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// TODO find User by user param
	// Set ctx
	c.Next()
}
