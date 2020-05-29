#!/bin/sh

set -e

replace_v2ray_config(){
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
                "clients":[
                    {
                        "id":"1e2c732f-4760-4e74-b654-9c4af7242b28",
                        "alterId":64,
                        "email":"noaway@gmail.com"
                    },
                    {
                        "id":"cefa7a9b-2f94-c7fd-b5a3-1e79d4ec8a15",
                        "alterId":10,
                        "email":"1014924101@qq.com"
                    },
                    {
                        "id":"bfb4aa24-8493-4986-9f49-f04ac9524adb",
                        "alterId":10,
                        "email":"maggie.hmg@hotmail.com"
                    },
                    {
                        "id":"6f27d2af-2bc6-42e9-a2b5-df4ab5049b85",
                        "alterId":10,
                        "email":"114408120@qq.com"
                    },
                    {
                        "id":"d89621cf-d80f-4449-af5f-7738d42ae44f",
                        "alterId":10,
                        "email":"1126866738@qq.com"
                    }
                ]
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
    exit 1
}
main