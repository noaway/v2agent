package cmd

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/spf13/cobra"
)

var noawayJob = `
job "noaway" {
	datacenters = ["conoha","taiwan"]
  
	  spread {
	  attribute = "${node.datacenter}"
	  weight = 100
	  }
  
	group "cache" {
	  task "redis" {
		driver = "docker"
  
		config {
		  image = "redis:3.2"
		  port_map {
			db = 6379
		  }
		}
  
		resources {
		  cpu    = 500
		  memory = 256
		  network {
			mbits = 10
			port "db" {}
		  }
		}
  
		service {
		  name = "redis-cache"
		  tags = ["global", "cache"]
		  port = "db"
		  check {
			name     = "alive"
			type     = "tcp"
			interval = "10s"
			timeout  = "2s"
		  }
		}
	  }
	}
  }  
`

func nomadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nomad",
		Short: "gen nomad",
		Long:  "generate a nomad",
		Run: func(_ *cobra.Command, _ []string) {
			config := api.DefaultConfig()
			client, err := api.NewClient(config)
			if err != nil {
				fmt.Println(err)
				return
			}
			jobs := client.Jobs()
			job, err := jobs.ParseHCL(noawayJob, false)
			if err != nil {
				fmt.Println("--- ", err)
				return
			}
			fmt.Println(*job)
			res, _, err := jobs.Register(job, nil)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(res.EvalID)
		},
	}
	return cmd
}
