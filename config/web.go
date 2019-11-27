package config

import(
	"github.com/noaway/godao"
)

type WebConfig struct{
	Addr string `hcl:"addr"`
}

func show(){
	_=godao.Engine
}
