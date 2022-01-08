package tasks

import (
	"encoding/json"
	"time"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/client/tasks/hardware"
	"github.com/kooiot/robot/pkg/net/msg"
)

type DoneTask struct {
	info    *msg.Task
	config  msg.DoneTask
	handler common.TaskHandler
	parent  common.Task
}

func init() {
	RegisterTask("usb", NewDoneTask)
}

func (s *DoneTask) Start() error {
	leds := []*hardware.NamedLed{}
	for _, name := range s.config.Leds {
		leds = append(leds, hardware.NewNamedLed(name))
	}
	// Make ourself as finished
	s.handler.OnSuccess(s)
	// Wait for other tasks
	s.handler.Wait(s.parent, func(task common.Task, result common.TaskResult) {
		if result.Result {
			// Halt device?
		} else {
			// Blink the LEDs
			for {
				for _, led := range leds {
					led.Set(255)
				}

				time.Sleep(200 * time.Millisecond)
				for _, led := range leds {
					led.Set(0)
				}

				time.Sleep(200 * time.Millisecond)
			}
		}
	})
	return nil
}

func (s *DoneTask) Stop() error {
	return nil
}

func NewDoneTask(handler common.TaskHandler, info *msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.DoneTask{}
	json.Unmarshal(data, &conf)

	return &DoneTask{
		info:    info,
		config:  conf,
		handler: handler,
		parent:  parent,
	}
}
