package tasks

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-ping/ping"
	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/log"
	uuid "github.com/satori/go.uuid"
)

type ModemTask struct {
	common.TaskBase
	config  msg.ModemTask
	handler common.TaskHandler
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

	log.Info("Modem task start ping")
	// pinger.SetPrivileged(true)
	// pinger.SetNetwork("ip4")
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		s.handler.OnError(s, err)
		return
	}

	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	stats_str, _ := json.Marshal(stats)

	if stats.PacketsRecv > 0 {
		opt := make(map[string]interface{})
		j, _ := json.Marshal(s.config.USB)
		json.Unmarshal(j, &opt)

		t := msg.Task{}
		t.UUID = uuid.NewV4().String()
		t.Name = "usb"
		t.Description = "Sub task from modem task"
		t.Option = opt

		s.handler.Spawn(NewUSBTask, &t, s)
		// s.handler.OnSuccess(s)
	} else {
		s.handler.OnError(s, errors.New("failed, statistics:"+string(stats_str)))
	}
}

func (s *ModemTask) Stop() error {
	return nil
}

func NewModemTask(handler common.TaskHandler, info *msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.ModemTask{}
	json.Unmarshal(data, &conf)

	return &ModemTask{
		TaskBase: common.NewTaskBase(info),
		config:   conf,
		handler:  handler,
	}
}
