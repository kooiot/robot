package tasks

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"time"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/client/tasks/hardware"
	"github.com/kooiot/robot/pkg/net/msg"
)

type USBTask struct {
	info    *msg.Task
	config  msg.USBTask
	handler common.TaskHandler
}

func init() {
	RegisterTask("usb", NewUSBTask)
}

func (s *USBTask) Start() error {
	go s.run()
	return nil
}

func (s *USBTask) findIds() (bool, error) {
	out, err := exec.Command("lsusb").Output()
	if err != nil {
		return false, err
	}

	found := true
	for i := 0; i < len(s.config.IDS); i++ {
		if !bytes.Contains(out, []byte(s.config.IDS[i])) {
			found = false
			err = errors.New("usb id " + s.config.IDS[i] + " missing")
		}
	}
	return found, err
}

func (s *USBTask) run() {
	time.Sleep(3 * time.Second)
	found, err := s.findIds()
	if !found {
		s.handler.OnError(s, err)
		return
	}
	if len(s.config.Reset) > 0 {
		// Test reset
		reset := hardware.NewNamedGPIO(s.config.Reset)

		reset.Set(1)
		time.Sleep(2 * time.Second)

		found, _ := s.findIds()
		if found {
			s.handler.OnError(s, errors.New("usb reset failed"))
			return
		}

		reset.Set(0)
		time.Sleep(2 * time.Second)

		found, _ = s.findIds()
		if !found {
			s.handler.OnError(s, errors.New("usb reset failed"))
			return
		}
	}
	if len(s.config.Power) > 0 {
		// Test reset
		power := hardware.NewNamedGPIO(s.config.Power)

		power.Set(0)
		time.Sleep(2 * time.Second)

		found, _ := s.findIds()
		if found {
			s.handler.OnError(s, errors.New("usb power failed"))
			return
		}
		power.Set(1)
		time.Sleep(2 * time.Second)

		found, _ = s.findIds()
		if !found {
			s.handler.OnError(s, errors.New("usb power failed"))
			return
		}
	}

	s.handler.OnSuccess(s)
}

func (s *USBTask) Stop() error {
	return nil
}

func NewUSBTask(handler common.TaskHandler, info *msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.USBTask{}
	json.Unmarshal(data, &conf)

	return &USBTask{
		info:    info,
		config:  conf,
		handler: handler,
	}
}
