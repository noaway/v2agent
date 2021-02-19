#!/bin/sh

set -e

replace_v2ray_config(){
    clients=`cat $(dirname $0)/clients.json`
    cat > /etc/v2ray/config.json << EOF
{
    "log":{
        "access":"/var/log/v2ray/access.log",
        "error":"/var/log/v2ray/error.log",
        "loglevel":"info"
    },
    "inbounds":[
        {
            "port":5219,
            "protocol":"vmess",
            "settings":{
                "clients": ${clients}
            },
            "streamSettings":{
                "network":"ws",
                "wsSettings":{
                    "path":"/echo"
                }
            },
            "tag":"proxy"
        },
        {
            "listen":"127.0.0.1",
            "port":9201,
            "protocol":"dokodemo-door",
            "settings":{
                "address":"127.0.0.1"
            },
            "tag":"api"
        }
    ],
    "outbounds":[
        {
            "protocol":"freedom",
            "settings":{

            }
        },
        {
            "protocol":"blackhole",
            "settings":{

            },
            "tag":"blocked"
        }
    ],
    "routing":{
        "rules":[
            {
                "inboundTag":[
                    "api"
                ],
                "outboundTag":"api",
                "type":"field"
            }
        ]
    },
    "stats":{

    },
    "api":{
        "services":[
            "StatsService",
            "HandlerService"
        ],
        "tag":"api"
    },
    "policy":{
        "levels":{
            "0":{
                "statsUserDownlink":true,
                "statsUserUplink":true
            }
        },
        "system":{
            "statsInboundUplink":true,
            "statsInboundDownlink":true
        }
    }
}
EOF
}

main(){
    replace_v2ray_config
    systemctl restart v2ray &&
    sleep 2 &&
    systemctl status v2ray
    exit 0
}
main