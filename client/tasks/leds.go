package tasks

import (
	"encoding/json"
	"time"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/client/tasks/hardware"
	"github.com/kooiot/robot/pkg/net/msg"
)

type LedsTask struct {
	common.TaskBase
	config  msg.LedsTask
	handler common.TaskHandler
	parent  common.Task
}

func init() {
	RegisterTask("usb", NewLedsTask)
}

func (s *LedsTask) Start() error {
	leds := []*hardware.NamedLed{}
	for _, name := range s.config.Leds {
		leds = append(leds, hardware.NewNamedLed(name))
	}

	for i := 0; i < s.config.Count; i++ {
		for _, led := range leds {
			led.Set(255)
		}
		time.Sleep(time.Millisecond * time.Duration(s.config.Span))
		for _, led := range leds {
			led.Set(0)
		}
		time.Sleep(time.Millisecond * time.Duration(s.config.Span))
	}
	// Make ourself as finished
	s.handler.OnSuccess(s)
	return nil
}

func (s *LedsTask) Stop() error {
	return nil
}

func NewLedsTask(handler common.TaskHandler, info *msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.LedsTask{}
	json.Unmarshal(data, &conf)

	return &LedsTask{
		TaskBase: common.NewTaskBase(info),
		config:   conf,
		handler:  handler,
		parent:   parent,
	}
}
