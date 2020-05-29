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
	V2ray map[string]config.V2CliConfig
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
		second := url.PathEscape(fmt.Sprintf("network=%v&wsPath=%v&aid=%v&tls=%v&allowInsecure=%v&remark=%v", v2ray.Protocol, v2ray.WSPath, v2ray.AlterId, tls, skipCertVerify, v2ray.Name))
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
