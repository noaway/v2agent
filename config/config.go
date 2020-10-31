package config

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Configuration struct {
	V2ray []V2CliConfig `hcl:"v2ray,block"`
	SS    []SSConfig    `hcl:"ss,block"`
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
