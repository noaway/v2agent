#!/bin/bash

set -e

nginx_conf_dir="/etc/nginx/conf.d"
host=""

is_root() {
    if [ 0 == $UID ]; then
        echo -e "当前用户是 root 用户"
    else
        echo -e "当前用户不是 root 用户, 请切换到 root 用户后重新执行脚本"
        exit 1
    fi
}

nginx_install(){
    if [ -x "$(command -V nginx)" ]; then
        echo -e "ngxin 已经安装过了, 跳过安装"
    else
        apt-get update &&
        apt-get install -y nginx
    fi
}

nginx_conf_add(){
    touch ${nginx_conf_dir}/server.conf
    cat > ${nginx_conf_dir}/server.conf <<EOF
server {
  listen 443 ssl;
  ssl on;
  ssl_certificate       /etc/letsencrypt/live/${host}/fullchain.pem;
  ssl_certificate_key   /etc/letsencrypt/live/${host}/privkey.pem;
  ssl_protocols         TLSv1 TLSv1.1 TLSv1.2;
  ssl_ciphers           HIGH:!aNULL:!MD5;
  server_name           conoha.noaway.io;
    location /echo {
      proxy_redirect off;
      proxy_pass http://127.0.0.1:5219;
      proxy_http_version 1.1;
      proxy_set_header Upgrade \$http_upgrade;
      proxy_set_header Connection "upgrade";
      proxy_set_header Host $host;
      # Show real IP in v2ray access.log
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
EOF
}

v2ray_install(){
    bash <(curl -L -s https://install.direct/go.sh)
}

certbot_install(){
    if ! [ -x "$(command -V certbot)" ]; then
        apt-get update &&
        apt-get install software-properties-common &&
        add-apt-repository universe &&
        add-apt-repository ppa:certbot/certbot &&
        apt-get update &&
        apt-get install certbot python-certbot-nginx &&
        certbot --nginx &&
        certbot certonly --nginx &&
        certbot renew --dry-run
    fi
}

open_bbr(){
    if ! [ -x "$(lsmod | grep bbr)" ]; then
        echo "net.core.default_qdisc=fq" >> /etc/sysctl.conf
        echo "net.ipv4.tcp_congestion_control=bbr" >> /etc/sysctl.conf
        sysctl -p
        sysctl net.ipv4.tcp_available_congestion_control
        lsmod | grep bbr
    else 
        echo -e "BBR 已设置"
        lsmod | grep bbr
    fi
}

nomad_setting(){
    read -rp "请输入DC: " dc
    if [ -z "$dc" ]; then
        echo "dc is empty"
        return
    fi
    read -rp "集群地址: " cluster
    if [ -z "$cluster" ]; then
        echo "cluster is empty"
        return
    fi
    read -rp "当前域名或 ip: " current_domain
    if [ -z "$current_domain" ]; then
        echo "current_domain is empty"
        return
    fi

    if ! [ -f "/usr/local/bin/nomad" ]; then
        echo "正在安装 nomad"
        wget https://releases.hashicorp.com/nomad/0.10.4/nomad_0.10.4_linux_amd64.zip &&
        unzip nomad_0.10.4_linux_amd64.zip
        chown root:root nomad
        mv nomad /usr/local/bin/
        echo "nomad version: $(nomad version)"
        nomad -autocomplete-install
        complete -C /usr/local/bin/nomad nomad
        mkdir --parents /opt/nomad
        touch /etc/systemd/system/nomad.service
        cat > /etc/systemd/system/nomad.service << EOF
[Unit]
Description=Nomad
Documentation=https://nomadproject.io/docs/
Wants=network-online.target
After=network-online.target

[Service]
ExecReload=/bin/kill -HUP $MAINPID
ExecStart=/usr/local/bin/nomad agent -config /etc/nomad.d
KillMode=process
KillSignal=SIGINT
LimitNOFILE=infinity
LimitNPROC=infinity
Restart=on-failure
RestartSec=2
StartLimitBurst=3
StartLimitIntervalSec=10
TasksMax=infinity

[Install]
WantedBy=multi-user.target
EOF
        mkdir --parents /etc/nomad.d
        chmod 700 /etc/nomad.d
        systemctl enable nomad
    fi

    touch /etc/nomad.d/nomad.hcl
    cat > /etc/nomad.d/nomad.hcl << EOF
datacenter = "$dc"
data_dir = "/opt/nomad"

server {
    enabled = true
    bootstrap_expect = 1
    server_join {
	    retry_join = ["$cluster"]
    }
}

client {
  enabled = true
  servers = ["127.0.0.1"]
}

advertise {
    rpc = "$current_domain"
    serf = "$current_domain"
}

plugin "raw_exec" {
  config {
    enabled = true
  }
}
EOF

    systemctl start nomad &&
    sleep 1 &&
    systemctl status nomad
}

main(){
    is_root

    while : 
    do
        echo -e "\t V2ray 管理脚本"

        echo -e "—————————————— 安装向导 ——————————————"""
        echo -e "1 安装 nginx"
        echo -e "2 添加默认 nginx 配置"
        echo -e "3 安装 v2ray"
        echo -e "4 安装 证书(安装前添加 DNS A 地址)"
        echo -e "5 配置 nomad"
        echo -e "—————————————— 查看信息 ——————————————"
        echo -e "6 重启 V2ray"
        echo -e "7 查看 V2ray 状态"
        echo -e "—————————————— 其他选项 ——————————————"
        echo -e "8 开启 BBR"
        echo -e "q 脚本退出"

        read -rp "请输入编号: " menu_num
        case $menu_num in
        1)
            nginx_install
            ;;
        2)
            nginx_conf_add
            ;;
        3)
            v2ray_install
            ;;
        4)
            certbot_install
            ;;
        5)
            nomad_setting
            ;;
        6)
            systemctl restart v2ray
            
            ;;
        7)
            systemctl status v2ray
            
            ;;
        8)
            open_bbr
            ;;
        q)
            echo -e "脚本退出"
            return
            ;;
        *)
            echo -e "请输入正确的数字"
            ;;
        esac
    done
}

main