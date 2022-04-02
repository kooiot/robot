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

type USBTask struct {
	common.TaskBase
	ctx     context.Context
	config  msg.USBTask
	handler common.TaskHandler
}

func init() {
	RegisterTask("usb", NewUSBTask)
}

func (s *USBTask) hasIds() (bool, error) {
	out, err := exec.Command("lsusb").Output()
	if err != nil {
		return false, err
	}
	// log.Info("lsusb output: %s %#v", out, s.config.IDS)

	ret := true
	for i := 0; i < len(s.config.IDS); i++ {
		if !bytes.Contains(out, []byte(s.config.IDS[i])) {
			ret = false
			err = errors.New("usb id " + s.config.IDS[i] + " missing")
			break
		}
	}
	return ret, err
}

func (s *USBTask) hasNoIds() (bool, error) {
	out, err := exec.Command("lsusb").Output()
	if err != nil {
		return false, err
	}
	// log.Info("lsusb output: %s", out)

	ret := true
	for i := 0; i < len(s.config.IDS); i++ {
		if bytes.Contains(out, []byte(s.config.IDS[i])) {
			ret = false
			err = errors.New("usb id " + s.config.IDS[i] + " found")
			break
		}
	}
	return ret, err
}

func (s *USBTask) Run() (interface{}, error) {
	xl := xlog.FromContextSafe(s.ctx)
	time.Sleep(3 * time.Second)
	ret, err := s.hasIds()
	if !ret {
		return nil, err
	}
	if len(s.config.Reset) > 0 {
		xl.Info("USB test reset gpio")
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
			ret, err = s.hasNoIds()
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
			ret, err = s.hasIds()
			if ret {
				break
			}
		}
		if !ret {
			return nil, errors.New("usb reset backover failed, error:" + err.Error())
		}
	}
	if len(s.config.Power) > 0 {
		xl.Info("USB test power gpio")
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
			ret, err = s.hasNoIds()
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
			ret, err = s.hasIds()
			if ret {
				break
			}
		}
		if !ret {
			return nil, errors.New("usb power up failed, error:" + err.Error())
		}
	}
	xl.Info("USB test done!")
	return "Done", nil
}

func NewUSBTask(ctx context.Context, handler common.TaskHandler, info msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.USBTask{}
	json.Unmarshal(data, &conf)
	// log.Info("USB Task: %#v from: %#v", conf, info.Option)
	xl := xlog.FromContextSafe(ctx).Spawn().AppendPrefix("Task.USB")

	return &USBTask{
		ctx:      xlog.NewContext(ctx, xl),
		TaskBase: common.NewTaskBase(info),
		config:   conf,
		handler:  handler,
	}
}
