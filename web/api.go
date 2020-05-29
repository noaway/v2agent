package web

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/noaway/v2agent/config"
	"github.com/noaway/v2agent/dispatch"
	"github.com/noaway/v2agent/internal/gensub"
	"github.com/noaway/v2agent/internal/utils"
	"github.com/noaway/v2agent/web/models"
	"github.com/sirupsen/logrus"
)

func addUser(dsp dispatch.DispatchHandle) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.PostForm("uuid")
		email := c.PostForm("email")
		region := c.PostForm("region")
		alterId := utils.StrTo(c.PostForm("alter_id")).MustUint32()
		conf := c.PostForm("conf")
		orderPackage := utils.StrTo(c.PostForm("package")).MustInt()
		price := utils.StrTo(c.PostForm("price")).MustInt64()

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
		now := time.Now()

		if err := models.SyncOrder(&models.Order{
			Email:      email,
			Package:    orderPackage,
			Config:     conf,
			Price:      price,
			ExpireTime: now.AddDate(0, orderPackage, 0),
		}); err != nil {
			c.JSON(200, gin.H{
				"errmsg": err,
			})
			return
		}

		if conf == "" {
			c.JSON(200, gin.H{
				"errmsg": "",
			})
			return
		}

		if config.Configure().SubscribePath == "" {
			c.JSON(200, gin.H{
				"errmsg": "SubscribePath is empty",
			})
			return
		}

		cli := struct {
			V2CliConfig map[string]config.V2CliConfig `hcl:"v2ray,block"`
		}{}

		if err := config.Unmarshal(email, []byte(conf), &cli); err != nil {
			c.JSON(200, gin.H{
				"errmsg": err,
			})
			return
		}

		paths := []string{}
	loop:
		for k, v := range gensub.KitMap {
			content := v.Content(gensub.ProxyConfig{V2ray: cli.V2CliConfig})
			relative := "/" + email + "/"
			path, err := utils.GetDir(fmt.Sprintf("%v/%v/", config.Configure().SubscribePath, email), utils.PathExists)
			if err != nil {
				c.JSON(200, gin.H{
					"errmsg": err,
				})
				return
			}
			switch k {
			case "quantumult":
				c := "quantumult.conf"
				path += c
				relative += c
			case "kitsunebi":
				c := "kitsunebi.txt"
				path += c
				relative += c
			default:
				continue loop
			}

			if err := utils.WriteFile(path, []byte(content)); err != nil {
				c.JSON(200, gin.H{
					"errmsg": err,
				})
				return
			}
			paths = append(paths, relative)
		}

		c.JSON(200, gin.H{
			"errmsg": "",
			"paths":  paths,
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
