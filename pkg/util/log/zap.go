package log

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var level zapcore.Level

func CreateLogger(log_dir string, filename string, log_level string, format string, encode_level string, log_in_console bool) (logger *zap.Logger) {
	fmt.Printf("dir %s log %s level %s\n", log_dir, filename, log_level)
	if ok, _ := PathExists(log_dir); !ok { // 判断是否有文件夹
		fmt.Printf("create %v directory\n", log_dir)
		_ = os.Mkdir(log_dir, os.ModePerm)
	}

	switch log_level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	// if level == zap.DebugLevel || level == zap.ErrorLevel {
	if level == zap.ErrorLevel {
		logger = zap.New(getEncoderCore(log_dir, filename, format, encode_level, log_in_console), zap.AddStacktrace(level))
	} else {
		logger = zap.New(getEncoderCore(log_dir, filename, format, encode_level, log_in_console))
	}
	// if config.TUN_CFG.Zap.ShowLine {
	// 	logger = logger.WithOptions(zap.AddCaller())
	// }
	return logger
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig(encode_level string) (ec zapcore.EncoderConfig) {
	ec = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case encode_level == "LowercaseLevelEncoder": // 小写编码器(默认)
		ec.EncodeLevel = zapcore.LowercaseLevelEncoder
	case encode_level == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		ec.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case encode_level == "CapitalLevelEncoder": // 大写编码器
		ec.EncodeLevel = zapcore.CapitalLevelEncoder
	case encode_level == "CapitalColorLevelEncoder": // 大写编码器带颜色
		ec.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		ec.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return ec
}

// getEncoder 获取zapcore.Encoder
func getEncoder(format string, encode_level string) zapcore.Encoder {
	if format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig(encode_level))
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig(encode_level))
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(log_dir string, filename string, format string, encode_level string, log_in_console bool) (core zapcore.Core) {
	writer, err := GetWriteSyncer(log_dir, filename, log_in_console) // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(getEncoder(format, encode_level), writer, level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05.000"))
}
