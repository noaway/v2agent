package config

import (
	"github.com/noaway/godao"
)

type WebConfig struct {
	Addr string                 `hcl:"addr"`
	DB   godao.PostgreSQLConfig `hcl:"db,block"`
	Salt string                 `htl:"salt"`
}
