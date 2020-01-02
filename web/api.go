package web

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/noaway/v2agent/config"
	"github.com/noaway/v2agent/dispatch"
	"github.com/noaway/v2agent/internal/gensub"
	"github.com/noaway/v2agent/internal/utils"
	"github.com/sirupsen/logrus"
)

func addUser(dsp dispatch.DispatchHandle) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.PostForm("uuid")
		email := c.PostForm("email")
		region := c.PostForm("region")
		alterId := utils.StrTo(c.PostForm("alter_id")).MustUint32()
		conf := c.PostForm("conf")

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
			V2CliConfig []config.V2CliConfig `hcl:"v2cli_config,block"`
		}{}

		if err := config.Unmarshal(email, []byte(conf), &cli); err != nil {
			c.JSON(200, gin.H{
				"errmsg": err,
			})
			return
		}
		cliConfigs := cli.V2CliConfig

		paths := []string{}
	loop:
		for k, v := range gensub.KitMap {
			content := v.Content(cliConfigs)
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
