job "sync_v2ray_config" {
  datacenters = ["taiwan","conoha","la"]
  type = "batch"

  # 分配到所有数据中心
  spread {
    attribute = "${node.datacenter}"
    weight    = 100
  }

  group "default" {
    # 有多少个数据中心就写几个
    count = 3

    task "v2ray_config" {
      driver = "raw_exec"
      config {
        command = "config.sh"
      }
      artifact {
        source = "https://conoha.noaway.io/subscribe/v2ray_config/config.sh"
      }
    }
  }
}