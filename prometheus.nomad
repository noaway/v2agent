job "prometheus" {
  datacenters = ["conoha"]
  type        = "service"

  group "monitoring" {
    count = 1

    restart {
      attempts = 10
      interval = "30m"
      delay    = "15s"
      mode     = "fail"
    }

    ephemeral_disk {
      size = 300
    }

    task "prometheus" {
      template {
        change_mode = "noop"
        destination = "local/prometheus.yml"

        data = <<EOH
---
global:
  scrape_interval:     5s
  evaluation_interval: 5s

scrape_configs:
  - job_name: 'nomad_metrics'
    static_configs:
    - targets: ['la.noaway.io:4646','conoha.noaway.io:4646','taiwan.noaway.io:4646']
    scrape_interval: 5s
    metrics_path: /v1/metrics
    params:
      format: ['prometheus']
  - job_name: 'v2ray_metrics'
    static_configs:
    - targets: ['conoha.noaway.io:9550','taiwan.noaway.io:9550','la.noaway.io:9550','hk.noaway.io:9550']
    scrape_interval: 5s
    metrics_path: /scrape
EOH
      }

      driver = "docker"

      config {
        image = "prom/prometheus:latest"
        network_mode = "host"

        volumes = [
          "local/prometheus.yml:/etc/prometheus/prometheus.yml",
        ]

        port_map {
          prometheus_ui = 9090
        }
      }

      resources {
        network {
          mbits = 10
          port "prometheus_ui" {
              static = 9090
          }
        }
      }
    }
  }
}