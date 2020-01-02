#!/bin/bash
set -e

installNginx(){
    if [ `command -v nginx` ];then
        echo "ok"
        return
    fi

    apt update && sudo apt install -y nginx
}

installV2ray(){
    bash <(curl -L -s https://install.direct/go.sh)
}

installCertbot(){
    apt-get update && apt-get install -y software-properties-common && add-apt-repository universe && add-apt-repository ppa:certbot/certbot && apt-get update
    apt-get install -y certbot python-certbot-nginx
}

settingBBR(){
    echo "net.core.default_qdisc=fq" >> /etc/sysctl.conf
    echo "net.ipv4.tcp_congestion_control=bbr" >> /etc/sysctl.conf
    sysctl -p
    net.ipv4.tcp_available_congestion_control = bbr cubic reno
    lsmod | grep bbr
}

main(){
    installNginx

    installV2ray

    installCertbot

    echo "ok"
}

main