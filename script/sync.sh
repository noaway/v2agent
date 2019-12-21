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

main(){
    installNginx

    echo "ok"
}

main