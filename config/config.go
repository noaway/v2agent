package config

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Configuration struct {
	// Log Log `hcl:"log,block"`
	// Server Server        `hcl:"server,block"`
	V2ray []V2CliConfig `hcl:"v2ray,block"`
	SS    []SSConfig    `hcl:"ss,block"`
}

type Log struct {
	LogLevel string `hcl:"log_level"`
	LogPath  string `hcl:"log_path"`
}

type Server struct {
	Addr string `hcl:"addr"`
}

var configure *Configuration

func Configure() *Configuration {
	if configure == nil {
		panic("Configuration no init")
	}
	return configure
}

func NewConfigure(filePath string) error {
	var config Configuration
	err := hclsimple.DecodeFile(filePath, nil, &config)
	if err != nil {
		return err
	}

	configure = &config
	return nil
}
