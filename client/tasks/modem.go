package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"time"

	"github.com/go-ping/ping"
	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/xlog"
	uuid "github.com/satori/go.uuid"
)

type ModemTask struct {
	common.TaskBase
	ctx     context.Context
	config  msg.ModemTask
	handler common.TaskHandler
}

func init() {
	RegisterTask("modem", NewModemTask)
}

func (s *ModemTask) Run() (interface{}, error) {
	xl := xlog.FromContextSafe(s.ctx)
	cmd := exec.Command("sh", "-c", "sysctl -w net.ipv4.ping_group_range=\"0   2147483647\"")
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	time.Sleep(5 * time.Second)
	pinger, err := ping.NewPinger(s.config.PingAddr)
	if err != nil {
		return nil, errors.New("pinger initialization failure")
	}
	pinger.Count = 3

	xl.Debug("Modem task start ping")
	// pinger.SetPrivileged(true)
	// pinger.SetNetwork("ip4")
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return nil, err
	}

	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	stats_str, _ := json.Marshal(stats)

	if stats.PacketsRecv > 0 {
		// opt := make(map[string]interface{})
		// j, _ := json.Marshal(s.config.USB)
		// json.Unmarshal(j, &opt)

		t := msg.Task{}
		t.UUID = uuid.NewV4().String()
		t.ID = s.TaskBase.ID() + ".usb"
		t.Task = "usb"
		t.Description = "Sub task from modem task"
		t.Option = s.config.USB

		s.handler.Spawn(NewUSBTask, t, s)
		return "Done", nil
	} else {
		return nil, errors.New("failed, statistics:" + string(stats_str))
	}
}

func NewModemTask(ctx context.Context, handler common.TaskHandler, info msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.ModemTask{}
	json.Unmarshal(data, &conf)

	xl := xlog.FromContextSafe(ctx).Spawn().AppendPrefix("Task.Modem")
	return &ModemTask{
		TaskBase: common.NewTaskBase(info),
		ctx:      xlog.NewContext(ctx, xl),
		config:   conf,
		handler:  handler,
	}
}
