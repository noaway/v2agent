package config

type V2CliConfig struct {
	Name           string `hcl:"name,label"`
	GroupName      string `hcl:"group_name"`
	Server         string `hcl:"server"`
	Port           int    `hcl:"port"`
	UUID           string `hcl:"uuid"`
	AlterId        int    `hcl:"alterId"`
	Cipher         string `hcl:"cipher"`
	Protocol       string `hcl:"protocol"`
	WSPath         string `hcl:"ws_path"`
	TLS            bool   `hcl:"tls"`
	TLSHost        string `hcl:"tls_host"`
	SkipCertVerify bool   `hcl:"skip_cert_verify"`
}
