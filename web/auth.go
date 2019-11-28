package web

import (
	"github.com/gin-gonic/gin"
)

func getLogin(c *gin.Context) {

}

func postLogin(c *gin.Context) {
	ctx := Ctx(c)
	username := c.Query("username")
	password := c.Query("password")
	if username == "admin" && password == "123456" {
		ctx.Session.Set(userkey, username)
		ctx.Session.Save()
		return
	}
	ctx.View(LOGIN)
}

func getRegister(c *gin.Context) {

}

func postRegister(c *gin.Context) {

}
