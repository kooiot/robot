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

func (t *RTCTask) Run() (interface{}, error) {
	xl := xlog.FromContextSafe(t.ctx)
	cmd := "hwclock -w"
	if len(t.config.File) > 0 {
		cmd += " -f " + t.config.File
	}
	xl.Debug("Run: %s", cmd)
	err := exec.Command("sh", "-c", cmd).Run()
	if err != nil {
		return nil, err
	}

	time.Sleep(10 * time.Second)
	cmd = "hwclock -r"
	if len(t.config.File) > 0 {
		cmd += " -f " + t.config.File
	}
	xl.Debug("Run: %s", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}
	time_len := len(time.ANSIC)
	rtc_now, _ := time.Parse(time.ANSIC, string(out[:time_len]))

	diff := time.Since(rtc_now)
	if diff > time.Second {
		return nil, errors.New("failed, diff:" + diff.String())
	} else {
		return "Done", nil
	}
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
