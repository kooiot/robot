package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"time"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type RTCTask struct {
	common.TaskBase
	ctx     context.Context
	config  msg.RTCTask
	handler common.TaskHandler
}

func init() {
	RegisterTask("rtc", NewRTCTask)
}

func (s *RTCTask) Start() error {
	go s.run()
	return nil
}

func (s *RTCTask) run() {
	cmd := "hwclock -w"
	if len(s.config.File) > 0 {
		cmd += " -f " + s.config.File
	}
	err := exec.Command("sh", "-c", cmd).Wait()
	if err != nil {
		s.handler.OnError(s, err)
		return
	}

	time.Sleep(10 * time.Second)
	cmd = "hwclock -r"
	if len(s.config.File) > 0 {
		cmd += " -f " + s.config.File
	}
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		s.handler.OnError(s, err)
		return
	}
	rtc_now, _ := time.Parse("2006-01-02 15:04:05", string(out[:19]))

	diff := time.Since(rtc_now)
	if diff > time.Second {
		s.handler.OnError(s, errors.New("failed, diff:"+diff.String()))
	} else {
		s.handler.OnSuccess(s)
	}
}

func (s *RTCTask) Stop() error {
	return nil
}

func NewRTCTask(ctx context.Context, handler common.TaskHandler, info msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.RTCTask{}
	json.Unmarshal(data, &conf)

	xl := xlog.FromContextSafe(ctx).Spawn().AppendPrefix("Task.RTC")
	return &RTCTask{
		TaskBase: common.NewTaskBase(info),
		ctx:      xlog.NewContext(ctx, xl),
		config:   conf,
		handler:  handler,
	}
}
