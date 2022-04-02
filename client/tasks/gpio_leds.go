package tasks

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Allenxuxu/toolkit/sync/atomic"
	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/client/tasks/hardware"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type GPIOLedsTask struct {
	common.TaskBase
	ctx     context.Context
	config  msg.LedsTask
	handler common.TaskHandler
	parent  common.Task
	stop    atomic.Bool
}

func init() {
	RegisterTask("gpio_leds", NewGPIOLedsTask)
}

func (t *GPIOLedsTask) Run() (interface{}, error) {
	gpios := []*hardware.NamedGPIO{}
	for _, name := range t.config.Leds {
		gpios = append(gpios, hardware.NewNamedGPIO(name))
	}

	for i := 0; i < t.config.Count; i++ {
		if t.stop.Get() {
			break
		}
		for _, gpio := range gpios {
			gpio.Set(255)
		}
		time.Sleep(time.Millisecond * time.Duration(t.config.Span))
		if t.stop.Get() {
			break
		}
		for _, led := range gpios {
			led.Set(0)
		}
		time.Sleep(time.Millisecond * time.Duration(t.config.Span))
	}
	return "Done", nil
}

func (t *GPIOLedsTask) Abort() error {
	t.stop.Set(true)
	return nil
}

func NewGPIOLedsTask(ctx context.Context, handler common.TaskHandler, info msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.LedsTask{}
	json.Unmarshal(data, &conf)

	xl := xlog.FromContextSafe(ctx).Spawn().AppendPrefix("Task.GPIO_LEDS")
	return &GPIOLedsTask{
		TaskBase: common.NewTaskBase(info),
		ctx:      xlog.NewContext(ctx, xl),
		config:   conf,
		handler:  handler,
		parent:   parent,
		stop:     atomic.Bool{},
	}
}
