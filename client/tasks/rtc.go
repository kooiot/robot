package tasks

import (
	"encoding/json"
	"os/exec"
	"time"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/pkg/net/msg"
)

type RTCTask struct {
	info    *msg.Task
	config  msg.RTCTask
	handler common.TaskHandler
}

func init() {
	RegisterTask("task", NewRTCTask)
}

func (s *RTCTask) Start() error {
	go s.run()
	return nil
}

func (s *RTCTask) run() {
	err := exec.Command("hwclock", "-w").Wait()
	if err != nil {
		s.handler.OnError(s, err)
		return
	}

	time.Sleep(10 * time.Second)
	out, err := exec.Command("hwclock", "-r").Output()
	if err != nil {
		s.handler.OnError(s, err)
		return
	}
	rtc_now, _ := time.Parse("2006-01-02 15:04:05", string(out[:19]))

	result := common.TaskResult{}
	diff := time.Since(rtc_now)
	if diff > time.Second {
		result.Result = false
		result.Error = "failed, diff:" + diff.String()
	} else {
		result.Result = true
		result.Error = "done!"
	}
	s.handler.OnResult(s.config, result)
}

func (s *RTCTask) Stop() error {
	return nil
}

func NewRTCTask(handler common.TaskHandler, info *msg.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.RTCTask{}
	json.Unmarshal(data, &conf)

	return &RTCTask{
		info:    info,
		config:  conf,
		handler: handler,
	}
}
