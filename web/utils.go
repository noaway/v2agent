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
	LOGIN    Page = "login.html"
	USER     Page = "user.html"
	REGISTER Page = "register.html"
)

func Ctx(c *gin.Context) *Context {
	value, exists := c.Get(DefaultKey)
	var ctx *Context
	if !exists || value == nil {
		ctx = &Context{Context: c}
	} else {
		ctx = value.(*Context)
	}
	ctx.Session = sessions.Default(c)
	if ctx.User == nil {
		ctx.User = new(User)
	}
	return ctx
}

type HandlerFunc func(*Context)
type Context struct {
	*gin.Context

	Session sessions.Session
	User    *User
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

func (c *Context) Redirect(location string) { c.Context.Redirect(http.StatusFound, location) }

type User struct {
	UserName string
}

func authRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		Ctx(c).Redirect("/auth/login")
		c.Abort()
		return
	}
	// TODO find User by user param
	// Set ctx
	c.Next()
}
