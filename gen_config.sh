#!/bin/bash

set -e

PROXY_PATH=/Users/noaway/Documents/proxy
URI=https://conoha.noaway.io/subscribe

run(){
    go run main.go conversion --kit $1 -c $2 > $1.conf
}

noExistsCreateDir(){
    echo "mkdir -p /home/noaway/subscribe/$1"
    if [[ $1 != "" ]]; then
        ssh noaway@conoha.noaway.io "mkdir -p /home/noaway/subscribe/$1"
    else
        echo "未知的参数 「$1」"
        exit 1
    fi
}

getSrc(){
    src=$PROXY_PATH/$1.hcl
    if [ -e $src ]; then
        echo $src
    else
        echo "未找到 「$src」 文件"
        exit 1
    fi
}

showUrl(){
    url=$1
    if [ `curl -s $url` ]; then
        echo $url
    else
        echo "failed"
    fi
}

uploadFile(){
    scp $1 noaway@conoha.noaway.io:/home/noaway/subscribe/$2
}

kitsunebi(){
    noExistsCreateDir $1
    local src=$(getSrc $1)
    run kitsunebi $src && uploadFile kitsunebi.conf $1 &&
    showUrl $URI/$1/kitsunebi.conf
}

default(){
    noExistsCreateDir $1
    local src=$(getSrc $1)
    run default $src && uploadFile default.conf $1 &&
    showUrl $URI/$1/default.conf
}

main(){
    case $1 in
    kitsunebi)
        kitsunebi $2
    ;;
    *)
        default $1
    ;;
    esac
}

main $@