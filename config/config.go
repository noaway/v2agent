package config

import (
	"fmt"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
)

type Configuration struct {
	Server          Server          `hcl:"server,block"`
	Node            Node            `hcl:"node,block"`
	V2rayAddr       string          `hcl:"v2ray_addr"`
	DataDir         string          `hcl:"data_dir"`
	V2HandlerConfig V2HandlerConfig `hcl:"v2ray_handler_service,block"`
	WebConfig       WebConfig       `hcl:"web_config,block"`
	V2CliConfig     []V2CliConfig   `hcl:"v2cli_config,block"`
}

type Server struct {
	HttpAddr string `hcl:"http_addr"`
}

type V2HandlerConfig struct {
	Addr string `hcl:"addr"`
	Tag  string `hcl:"tag"`
}

type Node struct {
	Name             string            `hcl:"name"`
	JoinClusterAddrs []string          `hcl:"join_cluster_addrs"`
	BindAddr         string            `hcl:"bind_addr"`
	AdvertiseAddr    string            `hcl:"advertise_addr"`
	Domain           string            `hcl:"domain"`
	Tag              map[string]string `hcl:"tag,block"`
}

var configure *Configuration

func Configure() *Configuration {
	if configure == nil {
		panic("Configuration no init")
	}
	return configure
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
}
