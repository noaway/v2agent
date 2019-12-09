package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/noaway/godao"
	"github.com/sirupsen/logrus"
)

type Configuration struct {
	Server          Server          `hcl:"server,block"`
	Agent           Agent           `hcl:"agent,block"`
	V2HandlerConfig V2HandlerConfig `hcl:"v2ray_handler_service,block"`
	V2CliConfig     []V2CliConfig   `hcl:"v2cli_config,block"`
	Log             Log             `hcl:"log,block"`
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
	DataDir          string   `hcl:"data_dir"`
	JoinClusterAddrs []string `hcl:"join_cluster_addrs"`
	BindAddr         string   `hcl:"bind_addr"`
	AdvertisePort    int      `hcl:"advertise_port"`
	Region           string   `hcl:"region"`
	// Tag              map[string]string `hcl:"tag,block"`
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
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(filename)

	if diags != nil {
		panic(fmt.Sprintf("NewConfigure.ParseHCLFile diags: %s", diags))
	}
	gohcl.DecodeBody(file.Body, nil, configure)

	Configure().Log.InitLogrus()
}
