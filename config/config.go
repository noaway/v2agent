package config

import (
	"fmt"
	"io/ioutil"
	"os"

	// "github.com/hashicorp/hcl2/gohcl"
	// "github.com/hashicorp/hcl2/hclparse"
	"github.com/hashicorp/hcl"
	"github.com/noaway/godao"
	"github.com/sirupsen/logrus"
)

type Configuration struct {
	Server          Server                 `hcl:"server,block"`
	Agent           Agent                  `hcl:"agent,block"`
	V2HandlerConfig V2HandlerConfig        `hcl:"v2ray_handler_service,block"`
	V2ray           map[string]V2CliConfig `hcl:"v2ray,block"`
	Ss              map[string]SsConfig    `hcl:"ss,block"`
	Log             Log                    `hcl:"log,block"`
	SubscribePath   string                 `hcl:"subscribe_path"`
}

// 本地服务器 http 配置
type Server struct {
	HttpAddr string                 `hcl:"http_addr"`
	DB       godao.PostgreSQLConfig `hcl:"db,block"`
	Salt     string                 `htl:"salt"`
}

// v2ray handler rpc 端口配置
type V2HandlerConfig struct {
	Addr string `hcl:"addr"`
	Tag  string `hcl:"tag"`
}

// 每个节点的配置
type Agent struct {
	Name             string   `hcl:"name"`
	SyncInterval     int      `hcl:"sync_interval"`
	DataDir          string   `hcl:"data_dir"`
	JoinClusterAddrs []string `hcl:"join_cluster_addrs"`
	BindAddr         string   `hcl:"bind_addr"`
	AdvertiseHost    string   `hcl:"advertise_host"` // 外网 host
	AdvertisePort    int      `hcl:"advertise_port"` // 外网 端口
	Region           string   `hcl:"region"`
}

type Log struct {
	Level   string   `hcl:"level"`
	LogPath string   `hcl:"log_path"`
	file    *os.File `hcl:"-"`
}

func (log *Log) InitLogrus() {
	if log.Level == "" {
		log.Level = "info"
	}

	level, err := logrus.ParseLevel(log.Level)
	if err != nil {
		panic("invalid log level")
	}
	logrus.SetLevel(level)
	log.LoggerToFile()
}

func (log *Log) LoggerToFile() {
	if log.LogPath == "" {
		return
	}

	file, err := os.Create(log.LogPath)
	if err != nil {
		panic("LoggerToFile err " + err.Error())
	}
	log.file = file
	logrus.SetOutput(file)
}

var configure *Configuration

func Configure() *Configuration {
	if configure == nil {
		panic("Configuration no init")
	}
	return configure
}

func Close() {
	file := Configure().Log.file
	if file != nil {
		file.Close()
	}
}

func NewConfigure(filename string) {
	configure = &Configuration{}
	if filename == "" {
		return
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("read config error. [err='%v',filename='%v']", err, filename))
	}

	if err := hcl.Decode(configure, string(data)); err != nil {
		panic(fmt.Sprintf("parse config file error. [err='%v',filename='%v']", err, filename))
	}

	Configure().Log.InitLogrus()
}

func Unmarshal(filename string, data []byte, v interface{}) error {
	return nil
}
