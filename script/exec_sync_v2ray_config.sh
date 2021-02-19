#!/bin/bash
set -e

main(){
    baseDir=$(dirname $0)/..
    cd $baseDir
    tar czvf config.tar.gz ./clients.json ./sync_v2ray_conf.sh
    scp $baseDir/config.tar.gz noaway@conoha.noaway.io:/home/noaway/subscribe/v2ray_config
}

main