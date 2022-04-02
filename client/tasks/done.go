package tasks

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/client/tasks/hardware"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type DoneTask struct {
	common.TaskBase
	ctx     context.Context
	xl      *xlog.Logger
	config  msg.DoneTask
	handler common.TaskHandler
	parent  common.Task
}

func init() {
	RegisterTask("done", NewDoneTask)
}

func (s *DoneTask) Blink() {
	leds := []*hardware.NamedLed{}
	gpios := []*hardware.NamedGPIO{}
	for _, name := range s.config.Leds {
		leds = append(leds, hardware.NewNamedLed(name))
	}
	for _, name := range s.config.GPIOLeds {
		gpios = append(gpios, hardware.NewNamedGPIO(name))
	}
	// Blink the LEDs
	for {
		for _, led := range leds {
			led.Set(255)
		}
		for _, gpio := range gpios {
			gpio.Set(1)
		}

		time.Sleep(200 * time.Millisecond)
		for _, led := range leds {
			led.Set(0)
		}
		for _, gpio := range gpios {
			gpio.Set(0)
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func (s *DoneTask) AllOn() {
	leds := []*hardware.NamedLed{}
	for _, name := range s.config.Leds {
		leds = append(leds, hardware.NewNamedLed(name))
	}
	// Blink the LEDs
	for {
		for _, led := range leds {
			led.Set(255)
		}
	}
}

func (s *DoneTask) Run() (interface{}, error) {
	xl := s.xl
	// Wait for other tasks
	err := s.handler.Wait(s.parent, func(task common.Task, result msg.TaskResultDetail) {
		xl.Info("done got result: %#v", result)
		if result.Result {
			r := s.handler.(*Runner)
			if r != nil && s.config.Halt {
				r.Halt()
			} else {
				s.AllOn()
			}
		} else {
			go s.Blink()
		}
	})
	if err != nil {
		return nil, err
	}
	return "Done", nil
}

func NewDoneTask(ctx context.Context, handler common.TaskHandler, info msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.DoneTask{}
	json.Unmarshal(data, &conf)

	xl := xlog.FromContextSafe(ctx).Spawn().AppendPrefix("Task.Done")
	return &DoneTask{
		TaskBase: common.NewTaskBase(info),
		ctx:      xlog.NewContext(ctx, xl),
		xl:       xl,
		config:   conf,
		handler:  handler,
		parent:   parent,
	}
}
