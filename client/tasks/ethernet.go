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
)

type EthernetTask struct {
	common.TaskBase
	ctx     context.Context
	config  msg.EthernetTask
	handler common.TaskHandler
}

func init() {
	RegisterTask("ethernet", NewEthernetTask)
}

func (s *EthernetTask) Init() error {
	cmd := exec.Command("sh", "-c", "sysctl -w net.ipv4.ping_group_range=\"0   2147483647\"")
	err := cmd.Run()
	if err != nil {
		return err
	}

	for _, v := range s.config.Init {
		cmd := exec.Command("sh", "-c", v)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *EthernetTask) Run() (interface{}, error) {
	xl := xlog.FromContextSafe(s.ctx)
	xl.Info("Ethernet task start")

	time.Sleep(3 * time.Second)
	pinger, err := ping.NewPinger(s.config.PingAddr)
	if err != nil {
		return nil, errors.New("pinger initialization failure")
	}
	pinger.Count = 3

	xl.Info("Ethernet task start ping")
	// pinger.SetPrivileged(true)
	// pinger.SetNetwork("ip4")
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return nil, err
	}

	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats

	xl.Info("Ethernet task ping Done")
	if stats.PacketsRecv > 0 {
		return stats, nil
	} else {
		stats_str, _ := json.Marshal(stats)
		return nil, errors.New("failed, statistics:" + string(stats_str))
	}
}

func NewEthernetTask(ctx context.Context, handler common.TaskHandler, info msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.EthernetTask{}
	json.Unmarshal(data, &conf)

	xl := xlog.FromContextSafe(ctx).Spawn().AppendPrefix("Task.Ethernet")
	return &EthernetTask{
		TaskBase: common.NewTaskBase(info),
		ctx:      xlog.NewContext(ctx, xl),
		config:   conf,
		handler:  handler,
	}
}
