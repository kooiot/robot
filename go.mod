module github.com/kooiot/robot

go 1.16

require (
	github.com/Allenxuxu/gev v0.3.1-0.20211110012922-7cee8af7cb57
	github.com/Allenxuxu/ringbuffer v0.0.11
	github.com/Allenxuxu/toolkit v0.0.1
	github.com/fsnotify/fsnotify v1.5.1
	github.com/go-ping/ping v0.0.0-20211014180314-6e2b003bffdd
	github.com/gobwas/pool v0.2.1
	github.com/npat-efault/crc16 v0.0.0-20161013170008-4128ccbe47c3
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	go.bug.st/serial v1.3.4
	go.uber.org/zap v1.19.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace github.com/Allenxuxu/gev => github.com/kooiot/gev v0.3.1-0.20220121071430-3e7c75272c11
