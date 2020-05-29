job "v2ray-exporter" {
  datacenters = ["taiwan","conoha","hk","la"]
  type = "service"

  # 分配到所有数据中心
  spread {
    attribute = "${node.datacenter}"
    weight    = 100
  }

  group "default" {
    # 有多少个数据中心就写几个
    count = 5

    restart {
      attempts = 10
      interval = "10m"
      delay    = "15s"
      mode     = "fail"
    }

    task "v2ray_config" {
      driver = "raw_exec"
      config {
        command = "v2ray-exporter_linux_amd64"
        args = ["--v2ray-endpoint","127.0.0.1:9201"]
      }
      artifact {
        source = "https://github.com/wi1dcard/v2ray-exporter/releases/latest/download/v2ray-exporter_linux_amd64"
      }
    }
  }
}