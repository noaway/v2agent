package config

import (
	"fmt"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclparse"
)

type Configuration struct {
	Name             string           `hcl:"name"`
	JoinClusterAddrs []string         `hcl:"join_cluster_addrs"`
	Server           Server           `hcl:"server,block"`
	V2rayAddr        string           `hcl:"v2ray_addr"`
	DataDir          string           `hcl:"data_dir"`
	V2HandlerService V2HandlerService `hcl:"v2ray_handler_service,block"`
}

type Server struct {
	HttpAddr      string `hcl:"http_addr"`
	BindAddr      string `hcl:"bind_addr"`
	AdvertiseAddr string `hcl:"advertise_addr"`
}

type V2HandlerService struct {
	Addr string `hcl:"addr"`
	Tag  string `hcl:"tag"`
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
