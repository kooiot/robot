package config

import (
	"fmt"

	"github.com/kooiot/robot/pkg/util/log"
)

type CommonConf struct {
	Addr     string `yaml:"addr" json:"addr"` // default 127.0.0.1
	Port     int    `yaml:"port" json:"port"`
	ClientID string `yaml:"client_id" json:"client_id"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	Hostname string `yaml:"hostname" json:"hostname"`
	Hardware string `yaml:"hardware" json:"hardware"`
	System   string `yaml:"system" json:"system"`
}

type ClientConf struct {
	Common CommonConf
	Log    log.LogConf
}

// GetDefaultClientConf returns a client configuration with default values.
func GetDefaultClientConf() ClientConf {
	return ClientConf{
		Common: CommonConf{
			Addr:     "127.0.0.1",
			Port:     7080,
			ClientID: "UNKNOWN",
			User:     "User",
			Password: "Password",
			Hostname: "Host",
			Hardware: "UNKNOWN",
			System:   "UNKNOWN",
		},
		Log: log.LogConf{
			Link:  "latest_log",
			Dir:   "log",
			Level: "info",
		},
	}
}

func ParseClientConfig(path string) (ClientConf, error) {
	cfg := GetDefaultClientConf()
	err := cfg.Load(path)
	return cfg, err
}

func (cfg *ClientConf) Complete() {
	// fmt.Printf("ProxyL: %v\n", cfg.Proxy)

	// if cfg.LogLink == "console" {
	// 	cfg.LogDir = "console"
	// } else {
	// 	cfg.LogDir = "file"
	// }
}

func (cfg *ClientConf) Validate() error {
	if cfg.Common.ClientID == "" || cfg.Common.ClientID == "UNKNOWN" {
		return fmt.Errorf("client id missing")
	}

	return nil
}

func (cfg *ClientConf) Load(path ...string) error {
	v := Viper(cfg, path...)
	if v == nil {
		return fmt.Errorf("invalid protocol")
	}
	cfg.Complete()

	return nil
}
