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

func (s *EthernetTask) Start() error {
	go s.run()
	return nil
}

func (s *EthernetTask) run() {
	xl := xlog.FromContextSafe(s.ctx)
	cmd := exec.Command("sh", "-c", "sysctl -w net.ipv4.ping_group_range=\"0   2147483647\"")
	err := cmd.Run()
	if err != nil {
		s.handler.OnError(s, err)
		return
	}

	for _, v := range s.config.Init {
		cmd := exec.Command("sh", "-c", v)
		err := cmd.Run()
		if err != nil {
			s.handler.OnError(s, err)
			return
		}
	}

	time.Sleep(3 * time.Second)
	pinger, err := ping.NewPinger(s.config.PingAddr)
	if err != nil {
		s.handler.OnError(s, errors.New("pinger initialization failure"))
		return
	}
	pinger.Count = 3

	xl.Debug("Ethernet task start ping")
	// pinger.SetPrivileged(true)
	// pinger.SetNetwork("ip4")
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		s.handler.OnError(s, err)
		return
	}

	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	stats_str, _ := json.Marshal(stats)

	if stats.PacketsRecv > 0 {
		s.handler.OnSuccess(s)
	} else {
		s.handler.OnError(s, errors.New("failed, statistics:"+string(stats_str)))
	}
}

func (s *EthernetTask) Stop() error {
	return nil
}

func NewEthernetTask(ctx context.Context, handler common.TaskHandler, info *msg.Task, parent common.Task) common.Task {
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
