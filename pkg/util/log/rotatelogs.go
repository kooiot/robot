package log

import (
	"os"
	"path"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap/zapcore"
)

func GetWriteSyncer(log_dir string, filename string, log_in_console bool) (zapcore.WriteSyncer, error) {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	log_filename := path.Join(log_dir, filename)
	pwd, err := os.Getwd()
	if !path.IsAbs(log_filename) {
		log_filename = path.Join(pwd, log_filename)
	}
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   log_filename,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     1, // days
		LocalTime:  true,
		Compress:   true,
	})

	if log_in_console {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}
