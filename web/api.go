package web

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/noaway/v2agent/dispatch"
	"github.com/noaway/v2agent/internal/utils"
	"github.com/sirupsen/logrus"
)

func addUser(dsp dispatch.DispatchHandle) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.PostForm("uuid")
		email := c.PostForm("email")
		region := c.PostForm("region")
		alterId := utils.StrTo(c.PostForm("alter_id")).MustUint32()

		if alterId == 0 {
			alterId = 64
		}

		u := &dispatch.User{}
		u.UUID = uuid
		u.Email = email
		u.AlterId = alterId
		u.Regions = strings.Split(region, ",")
		if err := dsp.AddUser(u); err != nil {
			c.JSON(200, gin.H{
				"errmsg": err,
			})
			return
		}
		c.JSON(200, gin.H{
			"errmsg": "",
		})
	}
}

func delUser(dsp dispatch.DispatchHandle) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")
		logrus.Debugf("del user api '%v'", email)
		if email == "" {
			c.JSON(200, gin.H{
				"errmsg": "email is empty",
			})
			return
		}
		if err := dsp.DelUser(email); err != nil {
			c.JSON(200, gin.H{
				"errmsg": err,
			})
			return
		}
		c.JSON(200, gin.H{
			"errmsg": "",
		})
	}
}
