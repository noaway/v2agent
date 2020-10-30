package gensub

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/noaway/v2agent/config"
	"github.com/noaway/v2agent/internal/utils"
)

var KitMap = map[string]Kit{
	"quantumult": NewQuantumult(),
	"kitsunebi":  NewKitsunebi(),
	"default":    NewV2rayDefaultKit(),
}

func encodeBase64(src string) string { return base64.RawStdEncoding.EncodeToString([]byte(src)) }

func format(f string, a ...interface{}) string { return fmt.Sprintf(f, a...) }

type ProxyConfig struct {
	V2ray []config.V2CliConfig
	Ss    map[string]config.SsConfig
}

type Kit interface {
	Content(ProxyConfig) string
	Subscribe() string
	URLSchema() string
}

func NewQuantumult() *Quantumult {
	return &Quantumult{}
}

type Quantumult struct {
	subscribe string
	urlSchema string
}

func (q *Quantumult) Content(proxy ProxyConfig) string {
	content := bytes.Buffer{}
	for _, v2ray := range proxy.V2ray {
		certificate := "0"
		if v2ray.SkipCertVerify {
			certificate = "1"
		}
		strs := []string{
			format("%v = vmess", v2ray.Name),
			v2ray.Server,
			utils.ToStr(v2ray.Port),
			v2ray.Cipher,
			format(`"%v"`, v2ray.UUID),
			format("group=%v", v2ray.GroupName),
			format("over-tls=%v", v2ray.TLS),
			format("tls-host=%v", v2ray.TLSHost),
			format("certificate=%v", certificate),
			format("obfs=%v", v2ray.Protocol),
			format(`obfs-path="%v"`, v2ray.WSPath),
			`obfs-header="Host: 01.alternate.19900101.xyz[Rr][Nn]User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 18_0_0 like Mac OS X) AppleWebKit/888.8.88 (KHTML, like Gecko) Mobile/6666666"`,
		}
		str := strings.Join(strs, ",")

		content.WriteString("vmess://" + encodeBase64(str))
		content.WriteString("\n")
	}
	return encodeBase64(content.String())
}

func (q *Quantumult) Subscribe() string { return q.subscribe }
func (q *Quantumult) URLSchema() string { return q.urlSchema }

/*
vmess://base64(security:uuid@host:port)?[key=urlencode(value)[&key=urlencode(value) ...]]


其中 base64、urlencode 为函数，security 为加密方式，最后一部分是以 & 为分隔符的参数列表，key 为参数名称，value 为相应的值，例如：network=kcp&aid=32&remark=服务器1 经过 urlencode 后为 network=kcp&aid=32&remark=%E6%9C%8D%E5%8A%A1%E5%99%A81


一个完整的例子：vmess://Y2hhY2hhMjAtcG9seTEzMDU6OTUxMzc4NTctNzBmYS00YWM4LThmOTAtNGUyMGFlYjY2MmNmQHVuaS5raXRzdW5lYmkuZnVuOjU2NjY=?network=ws&wsPath=/v2&aid=0&tls=1&allowInsecure=1&mux=0&muxConcurrency=8&remark=WSS%20Test%20Outbound


可选参数（参数名称不区分大小写）：

network - 可选的值为 "tcp"、 "kcp"、"ws"、"h2" 等

wsPath - WebSocket 的协议路径

wsHost - WebSocket HTTP 头里面的 Host 字段值

kcpHeader - kcp 的伪装类型

uplinkCapacity - kcp 的上行容量

downlinkCapacity - kcp 的下行容量

h2Path - h2 的路径

h2Host - h2 的域名

quicSecurity - quic 加密方式

quicKey - quic 加密密钥

quicHeaderType - quic 头部伪装类型

aid - AlterId

tls - 是否启用 TLS，为 0 或 1

allowInsecure - TLS 的 AllowInsecure，为 0 或 1

tlsServer - TLS 的服务器端证书的域名

mux - 是否启用 mux，为 0 或 1

muxConcurrency - mux 的 最大并发连接数

remark - 备注名称


导入配置时，不在列表中的参数一般会按照 Core 的默认值处理。


ss:// 和 socks:// 的格式类似。
*/

func NewKitsunebi() *Kitsunebi {
	return &Kitsunebi{}
}

type Kitsunebi struct {
	Host  string `json:"host"`
	Path  string `json:"path"`
	Tls   string `json:"tls"`
	Add   string `json:"add"`
	Port  int    `json:"port"`
	Aid   int    `json:"aid"`
	Net   string `json:"net"`
	Type  string `json:"type"`
	V     string `json:"v"`
	PS    string `json:"ps"`
	ID    string `json:"id"`
	Class int    `json:"class"`
}

func (kit *Kitsunebi) Content(proxy ProxyConfig) string {
	content := bytes.Buffer{}
	for _, v2ray := range proxy.V2ray {
		first := encodeBase64(fmt.Sprintf("%v:%v@%v:%v", v2ray.Cipher, v2ray.UUID, v2ray.Server, v2ray.Port))
		tls := 1
		if !v2ray.TLS {
			tls = 0
		}
		skipCertVerify := 1
		if !v2ray.SkipCertVerify {
			skipCertVerify = 0
		}
		second := url.PathEscape(fmt.Sprintf("network=%v&wsPath=%v&aid=%v&tls=%v&allowInsecure=%v&remark=%v", v2ray.Protocol, v2ray.WSPath, v2ray.AlterId, tls, skipCertVerify, url.QueryEscape(v2ray.Name)))
		content.WriteString("vmess://" + first + "?" + second)
		content.WriteString("\n")
	}

	return encodeBase64(content.String())
}

func (kit *Kitsunebi) Subscribe() string { return "" }

func (kit *Kitsunebi) URLSchema() string { return "" }

func NewV2rayDefaultKit() *V2rayDefaultKit {
	return &V2rayDefaultKit{}
}

type V2rayDefaultKit struct {
	Host  string `json:"host"`
	Path  string `json:"path"`
	Tls   string `json:"tls"`
	Add   string `json:"add"`
	Port  int    `json:"port"`
	Aid   int    `json:"aid"`
	Net   string `json:"net"`
	Type  string `json:"type"`
	V     string `json:"v"`
	PS    string `json:"ps"`
	ID    string `json:"id"`
	Class int    `json:"class"`
}

func (kit *V2rayDefaultKit) Content(proxy ProxyConfig) string {
	content := bytes.Buffer{}
	for _, v2ray := range proxy.V2ray {
		kit.Host = v2ray.Server
		kit.Path = v2ray.WSPath
		if v2ray.TLS {
			kit.Tls = "tls"
		}
		kit.Add = v2ray.Server
		kit.Port = v2ray.Port
		kit.Aid = v2ray.AlterId
		kit.Net = v2ray.Protocol
		kit.Type = v2ray.Cipher
		kit.V = "2"
		kit.PS = v2ray.Name
		kit.ID = v2ray.UUID
		kit.Class = 0

		data, err := json.Marshal(kit)
		if err != nil {
			fmt.Println("Kitsunebi.Marshal err, ", err)
			return ""
		}

		str := string(data)
		content.WriteString("vmess://" + encodeBase64(str))
		content.WriteString("\n")
	}

	return encodeBase64(content.String())
}

func (kit *V2rayDefaultKit) Subscribe() string { return "" }

func (kit *V2rayDefaultKit) URLSchema() string { return "" }
