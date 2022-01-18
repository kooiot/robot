package log

import (
	"fmt"

	"go.uber.org/zap"
)

// Log is the under log object
var Log *zap.Logger

type LogConf struct {
	Link  string `yaml:"link" json:"link"`
	Dir   string `yaml:"dir" json:"dir"`
	Level string `yaml:"level" json:"level"`
}

func init() {
}

func InitLog(conf LogConf) {
	Log = CreateLogger(conf.Dir, conf.Link, conf.Level, "console", "LowercaseColorLevelEncoder", true)
}

// wrap log

func Error(format string, v ...interface{}) {
	Log.Error(fmt.Sprintf(format, v...))
}

func Warn(format string, v ...interface{}) {
	Log.Warn(fmt.Sprintf(format, v...))
}

func Info(format string, v ...interface{}) {
	Log.Info(fmt.Sprintf(format, v...))
}

func Debug(format string, v ...interface{}) {
	Log.Debug(fmt.Sprintf(format, v...))
}

func Trace(format string, v ...interface{}) {
	Log.Debug(fmt.Sprintf(format, v...))
}
