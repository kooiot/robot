package tasks

import (
	"encoding/json"
	"time"

	"github.com/go-ping/ping"
	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/pkg/net/msg"
	uuid "github.com/satori/go.uuid"
)

type ModemTask struct {
	handler common.TaskHandler
	config  msg.ModemTask
}

func init() {
	RegisterTask("modem", NewModemTask)
}

func (s *ModemTask) Start() error {
	go s.run()
	return nil
}

func (s *ModemTask) run() {
	time.Sleep(5 * time.Second)
	pinger, err := ping.NewPinger(s.config.PingAddr)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		panic(err)
	}

	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	stats_str, _ := json.Marshal(stats)

	result := common.TaskResult{}
	if stats.PacketsRecv > 0 {

		t := msg.Task{}
		t.UUID = uuid.NewV4().String()
		t.Name = "usb"
		t.Description = "Sub task from modem task"
		t.Option = s.config.USB

		s.handler.Spawn(NewUSBTask, &t, s)

		result.Result = true
		result.Error = string(stats_str)
		s.handler.OnResult(s.config, result)
	} else {
		result.Result = false
		result.Error = "failed, statistics:" + string(stats_str)
		s.handler.OnResult(s.config, result)
	}
}

func (s *ModemTask) Stop() error {
	return nil
}

func NewModemTask(handler common.TaskHandler, option interface{}) common.Task {
	data, _ := json.Marshal(option)

	conf := msg.ModemTask{}
	json.Unmarshal(data, &conf)

	return &ModemTask{
		handler: handler,
		config:  conf,
	}
}
