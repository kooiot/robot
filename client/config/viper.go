package config

import (
	"flag"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Viper(config interface{}, path ...string) *viper.Viper {
	var config_file string
	if len(path) == 0 {
		flag.StringVar(&config_file, "c", "", "choose config file.")
		flag.Parse()
		if config_file == "" { // 优先级: 命令行 > 环境变量 > 默认值
			config_file = "config.yaml"
			fmt.Printf("您正在使用config的默认值,config的路径为%v\n", config_file)
		} else {
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%v\n", config_file)
		}
	} else {
		config_file = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%v\n", config_file)
	}

	v := viper.New()
	v.SetConfigFile(config_file)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&config); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&config); err != nil {
		fmt.Println(err)
		return nil
	}

	return v
}
