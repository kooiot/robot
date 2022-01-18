package log

import (
	"os"

	"go.uber.org/zap"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			Log.Debug("create directory" + v)
			err = os.MkdirAll(v, os.ModePerm)
			if err != nil {
				Log.Error("create directory"+v, zap.Any(" error:", err))
			}
		}
	}
	return err
}
