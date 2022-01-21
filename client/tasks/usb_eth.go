package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"time"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/client/tasks/hardware"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type USBETHTask struct {
	common.TaskBase
	ctx     context.Context
	config  msg.USBETHTask
	handler common.TaskHandler
}

func init() {
	RegisterTask("usb_eth", NewUSBTask)
}

func (s *USBETHTask) hasEth() (bool, error) {
	out, err := exec.Command("sh", "-c", "cat /proc/net/dev").Output()
	if err != nil {
		return false, err
	}

	if !bytes.Contains(out, []byte(s.config.Name+":")) {
		return false, errors.New("ethernet " + s.config.Name + " missing")
	}
	return true, nil
}

func (s *USBETHTask) Run() (interface{}, error) {
	time.Sleep(3 * time.Second)
	ret, err := s.hasEth()
	if !ret {
		return nil, err
	}
	if len(s.config.Reset) > 0 {
		// Test reset
		reset := hardware.NewNamedGPIO(s.config.Reset)

		// pull reset to high
		err = reset.Set(1)
		if nil != err {
			return nil, err
		}
		now := time.Now()
		var ret bool
		for time.Since(now) < 5*time.Second {
			time.Sleep(500 * time.Millisecond)
			ret, err = s.hasEth()
			if ret {
				break
			}
		}
		if !ret {
			return nil, errors.New("usb reset trigger failed, error:" + err.Error())
		}

		// pull reset to low
		err = reset.Set(0)
		if nil != err {
			return nil, err
		}

		now = time.Now()
		for time.Since(now) < 15*time.Second {
			time.Sleep(500 * time.Millisecond)
			ret, err = s.hasEth()
			if ret {
				break
			}
		}
		if !ret {
			return nil, errors.New("usb reset backover failed, error:" + err.Error())
		}
	}
	if len(s.config.Power) > 0 {
		// Test reset
		power := hardware.NewNamedGPIO(s.config.Power)

		err = power.Set(0)
		if nil != err {
			return nil, err
		}
		now := time.Now()
		var ret bool
		for time.Since(now) < 5*time.Second {
			time.Sleep(500 * time.Millisecond)
			ret, err = s.hasEth()
			if ret {
				break
			}
		}
		if !ret {
			return nil, errors.New("usb power down failed, error:" + err.Error())
		}
		err = power.Set(1)
		if nil != err {
			return nil, err
		}

		now = time.Now()
		for time.Since(now) < 15*time.Second {
			time.Sleep(500 * time.Millisecond)
			ret, err = s.hasEth()
			if ret {
				break
			}
		}
		if !ret {
			return nil, errors.New("usb power up failed, error:" + err.Error())
		}
	}
	return "Done", nil
}

func NewUSBETHTask(ctx context.Context, handler common.TaskHandler, info msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.USBETHTask{}
	json.Unmarshal(data, &conf)
	// log.Info("USB Task: %#v from: %#v", conf, info.Option)
	xl := xlog.FromContextSafe(ctx).Spawn().AppendPrefix("Task.USB")

	return &USBETHTask{
		ctx:      xlog.NewContext(ctx, xl),
		TaskBase: common.NewTaskBase(info),
		config:   conf,
		handler:  handler,
	}
}
