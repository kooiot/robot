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
	handler common.TaskHandler
	config  msg.USBTask
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
		if bytes.Index(out, []byte(s.config.IDS[i])) < 0 {
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
		result := common.TaskResult{}
		result.Result = false
		result.Error = err.Error()
		s.handler.OnResult(s.config, result)
		return
	}
	if len(s.config.Reset) > 0 {
		// Test reset
		reset := hardware.NewNamedGPIO(s.config.Reset)
		reset.Set(1)
		time.Sleep(2 * time.Second)

		found, _ := s.findIds()
		if found {
			result := common.TaskResult{}
			result.Result = false
			result.Error = "usb reset failed"
			s.handler.OnResult(s.config, result)
			return
		}
		reset.Set(0)
	}
	if len(s.config.Reset) > 0 {
		// Test reset
		power := hardware.NewNamedGPIO(s.config.Power)
		power.Set(0)
		time.Sleep(2 * time.Second)

		found, _ := s.findIds()
		if found {
			result := common.TaskResult{}
			result.Result = false
			result.Error = "usb reset failed"
			s.handler.OnResult(s.config, result)
			return
		}
		power.Set(1)
	}

	result := common.TaskResult{}
	result.Result = true
	result.Error = "done!"
	s.handler.OnResult(s.config, result)
}

func (s *USBTask) Stop() error {
	return nil
}

func NewUSBTask(handler common.TaskHandler, option interface{}) common.Task {
	data, _ := json.Marshal(option)

	conf := msg.USBTask{}
	json.Unmarshal(data, &conf)

	return &USBTask{
		handler: handler,
		config:  conf,
	}
}
