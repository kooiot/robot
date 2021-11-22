package config

import (
	"fmt"

	"github.com/kooiot/robot/pkg/util/log"
)

type CommonConf struct {
	Bind  string `mapstructure:"bind" json:"bind"`   // default 0.0.0.0
	Port  int    `mapstructure:"port" json:"port"`   // default 7080
	Loops int    `mapstructure:"loops" json:"loops"` // default is 0
}

type AutoMatch struct {
	Key   string `mapstructure:"key" json:"key"`
	Match string `mapstructure:"match" json:"match"`
}

type AutoConfig struct {
	Matches []AutoMatch `mapstructure:"matches" json:"matches"`
	Config  string      `mapstructure:"config" json:"config"`
}

type AutoTasks struct {
	Folder string       `mapstructure:"folder" json:"folder"`
	Autos  []AutoConfig `mapstructure:"autos" json:"autos"`
}

type ServerConf struct {
	Common CommonConf  `mapstructure:"common" json:"common"`
	Log    log.LogConf `mapstructure:"log" json:"log"`
	Tasks  AutoTasks   `mapstructure:"tasks" json:"tasks"`
}

// GetDefaultServerConf returns a client configuration with default values.
func GetDefaultServerConf() ServerConf {
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

func ParseServerConfig(path string) (ServerConf, error) {
	cfg := GetDefaultServerConf()
	err := cfg.Load(path)
	return cfg, err
}

func (cfg *ServerConf) Complete() {
	// fmt.Printf("Tasks: %v\n", cfg.Tasks)

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
