package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func getLogin(c *Context) {
	c.View(LOGIN, gin.H{
		"Action":      "/auth/login",
		"RegisterURL": "/auth/register",
	})
}

func postLogin(c *Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	fmt.Println("email: ", email, " password: ", password)
	if email == "admin@gmail.com" && password == "123456" {
		c.Session.Set(userkey, email)
		c.Session.Save()
		c.Redirect("/user")
		return
	}
	c.View(LOGIN)
}

func getRegister(c *Context) {
	c.View(REGISTER, gin.H{
		"Action":   "/auth/register",
		"LoginURL": "/auth/login",
	})
}

func postRegister(c *Context) {

}
