package config

import (
	"fmt"

	"github.com/kooiot/robot/pkg/util/log"
)

type CommonConf struct {
	Bind  string `yaml:"bind" json:"bind"` // default 0.0.0.0
	Port  int    `yaml:"port" json:"port"`
	Loops int    `yaml:"loops" json:"loops"`
}

type ServerConf struct {
	Common CommonConf
	Log    log.LogConf
}

// GetDefaultClientConf returns a client configuration with default values.
func GetDefaultClientConf() ServerConf {
	return ServerConf{
		Common: CommonConf{
			Bind: "0.0.0.0",
			Port: 7080,
		},
		Log: log.LogConf{
			Link:  "latest_log",
			Dir:   "log",
			Level: "info",
		},
	}
}

func ParseClientConfig(path string) (ServerConf, error) {
	cfg := GetDefaultClientConf()
	err := cfg.Load(path)
	return cfg, err
}

func (cfg *ServerConf) Complete() {
	// fmt.Printf("ProxyL: %v\n", cfg.Proxy)

	// if cfg.LogLink == "console" {
	// 	cfg.LogDir = "console"
	// } else {
	// 	cfg.LogDir = "file"
	// }
}

func (cfg *ServerConf) Validate() error {
	// if cfg.Common.ClientID == "" || cfg.Common.ClientID == "UNKNOWN" {
	// 	return fmt.Errorf("client id missing")
	// }

	return nil
}

func (cfg *ServerConf) Load(path ...string) error {
	v := Viper(cfg, path...)
	if v == nil {
		return fmt.Errorf("invalid protocol")
	}
	cfg.Complete()

	return nil
}
