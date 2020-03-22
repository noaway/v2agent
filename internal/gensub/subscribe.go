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

type Kit interface {
	Content([]config.V2CliConfig) string
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

func (q *Quantumult) Content(confs []config.V2CliConfig) string {
	content := bytes.Buffer{}
	for i, conf := range confs {
		certificate := "0"
		if conf.SkipCertVerify {
			certificate = "1"
		}
		strs := []string{
			format("%v = vmess", conf.Name),
			conf.Server,
			utils.ToStr(conf.Port),
			conf.Cipher,
			format(`"%v"`, conf.UUID),
			format("group=%v", conf.GroupName),
			format("over-tls=%v", conf.TLS),
			format("tls-host=%v", conf.TLSHost),
			format("certificate=%v", certificate),
			format("obfs=%v", conf.Protocol),
			format(`obfs-path="%v"`, conf.WSPath),
			`obfs-header="Host: 01.alternate.19900101.xyz[Rr][Nn]User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 18_0_0 like Mac OS X) AppleWebKit/888.8.88 (KHTML, like Gecko) Mobile/6666666"`,
		}
		str := strings.Join(strs, ",")

		content.WriteString("vmess://" + encodeBase64(str))
		if i < len(confs)-1 {
			content.WriteString("\n")
		}
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

func (kit *Kitsunebi) Content(confs []config.V2CliConfig) string {
	content := bytes.Buffer{}
	for i, conf := range confs {
		first := encodeBase64(fmt.Sprintf("%v:%v@%v:%v", conf.Cipher, conf.UUID, conf.Server, conf.Port))
		tls := 1
		if !conf.TLS {
			tls = 0
		}
		skipCertVerify := 1
		if !conf.SkipCertVerify {
			skipCertVerify = 0
		}
		second := url.PathEscape(fmt.Sprintf("network=%v&wsPath=%v&aid=%v&tls=%v&allowInsecure=%v&remark=%v", conf.Protocol, conf.WSPath, conf.AlterId, tls, skipCertVerify, conf.Name))
		content.WriteString("vmess://" + first + "?" + second)

		if i < len(confs)-1 {
			content.WriteString("\n")
		}
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

func (kit *V2rayDefaultKit) Content(confs []config.V2CliConfig) string {
	content := bytes.Buffer{}
	for i, conf := range confs {
		kit.Host = conf.Server
		kit.Path = conf.WSPath
		if conf.TLS {
			kit.Tls = "tls"
		}
		kit.Add = conf.Server
		kit.Port = conf.Port
		kit.Aid = conf.AlterId
		kit.Net = conf.Protocol
		kit.Type = conf.Cipher
		kit.V = "2"
		kit.PS = conf.Name
		kit.ID = conf.UUID
		kit.Class = 0

		data, err := json.Marshal(kit)
		if err != nil {
			fmt.Println("Kitsunebi.Marshal err, ", err)
			return ""
		}

		str := string(data)
		content.WriteString("vmess://" + encodeBase64(str))
		if i < len(conf.Server)-1 {
			content.WriteString("\n")
		}
	}

	return encodeBase64(content.String())
}

func (kit *V2rayDefaultKit) Subscribe() string { return "" }

func (kit *V2rayDefaultKit) URLSchema() string { return "" }
