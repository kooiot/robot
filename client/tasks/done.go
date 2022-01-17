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
	config  msg.DoneTask
	handler common.TaskHandler
	parent  common.Task
}

func init() {
	RegisterTask("done", NewDoneTask)
}

func (s *DoneTask) Blink() {
	leds := []*hardware.NamedLed{}
	for _, name := range s.config.Leds {
		leds = append(leds, hardware.NewNamedLed(name))
	}
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

func (s *DoneTask) Start() error {
	// Make ourself as finished
	s.handler.OnSuccess(s)
	// Wait for other tasks
	s.handler.Wait(s.parent, func(task common.Task, result msg.TaskResult) {
		if result.Result {
			r := s.handler.(*Runner)
			if r != nil {
				r.Halt()
			}
		} else {
			go s.Blink()
		}
	})
	return nil
}

func (s *DoneTask) Stop() error {
	return nil
}

func NewDoneTask(ctx context.Context, handler common.TaskHandler, info msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.DoneTask{}
	json.Unmarshal(data, &conf)

	xl := xlog.FromContextSafe(ctx).Spawn().AppendPrefix("Task.Done")
	return &DoneTask{
		TaskBase: common.NewTaskBase(info),
		ctx:      xlog.NewContext(ctx, xl),
		config:   conf,
		handler:  handler,
		parent:   parent,
	}
}
