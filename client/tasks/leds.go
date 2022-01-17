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

type LedsTask struct {
	common.TaskBase
	ctx     context.Context
	config  msg.LedsTask
	handler common.TaskHandler
	parent  common.Task
}

func init() {
	RegisterTask("leds", NewLedsTask)
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

func NewLedsTask(ctx context.Context, handler common.TaskHandler, info msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.LedsTask{}
	json.Unmarshal(data, &conf)

	xl := xlog.FromContextSafe(ctx).Spawn().AppendPrefix("Task.LEDS")
	return &LedsTask{
		TaskBase: common.NewTaskBase(info),
		ctx:      xlog.NewContext(ctx, xl),
		config:   conf,
		handler:  handler,
		parent:   parent,
	}
}
