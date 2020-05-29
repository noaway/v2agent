job "grafana" {
  datacenters = ["conoha"]
  type = "service"

  group "monitoring" {
    restart {
      attempts = 10
      interval = "5m"
      delay = "10s"
      mode = "delay"
    }

    volume "grafana" {
      type      = "host"
      read_only = false
      source    = "grafana"
    }

    task "grafana" {
      driver = "docker"
      
      volume_mount {
        volume      = "grafana"
        destination = "/etc/grafana"
        read_only   = false
      }

      config {
        image = "grafana/grafana"
        network_mode = "host" 

        port_map = {
          http = 3000
        }
      }

      resources {
        cpu    = 128
        memory = 128
        network {
          mbits = 10
          port "http" {
              static = 3000
          }
        }
      }
    }
  }
}