package tasks

import (
	"context"
	"encoding/json"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/client/helper"
	"github.com/kooiot/robot/client/port"
	"github.com/kooiot/robot/client/port/serial"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type SerialResult struct {
	Ping common.StreamResult `json:"ping"`
	Pong common.StreamResult `json:"pong"`
}

type SerialTask struct {
	common.TaskBase
	ctx     context.Context
	config  msg.SerialTask
	handler common.TaskHandler

	src_port *serial.SerialPort
	src      *helper.PingPong
	dst_port *serial.SerialPort
	dst      *helper.PingPong
}

func init() {
	RegisterTask("serial", NewSerialTask)
}

func (s *SerialTask) Run() (interface{}, error) {
	err := s.src_port.Open()
	if err != nil {
		return nil, err
	}
	err = s.dst_port.Open()
	if err != nil {
		return nil, err
	}
	s.src.Start()
	s.dst.Start()
	s.src.Wait()
	s.dst.Wait()

	return &msg.TaskResultDetail{
		Result: true,
		Info:   "Done",
		Detail: SerialResult{
			Ping: s.src.Result,
			Pong: s.src.Result,
		},
	}, nil
}

func (s *SerialTask) Abort() error {
	s.src.Stop()
	s.dst.Stop()
	return nil
}

func NewSerialTask(ctx context.Context, handler common.TaskHandler, info msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.SerialTask{}
	json.Unmarshal(data, &conf)

	src_stream := port.NewStream()
	src, err := serial.NewSerial(src_stream, serial.Port(conf.SrcPort), serial.Baudrate(conf.Baudrate))
	if err != nil {
		return nil
	}

	dest_stream := port.NewStream()
	dest, err := serial.NewSerial(dest_stream, serial.Port(conf.DestPort), serial.Baudrate(conf.Baudrate))
	if err != nil {
		return nil
	}
	src_config := helper.PingPongConfig{IsPing: true, Count: conf.Count, MaxMsgSize: conf.MaxMsgSize}
	dest_config := helper.PingPongConfig{IsPing: false, Count: conf.Count, MaxMsgSize: conf.MaxMsgSize}

	xl := xlog.FromContextSafe(ctx).Spawn().AppendPrefix("Task.Serial")

	t := &SerialTask{
		TaskBase: common.NewTaskBase(info),
		ctx:      xlog.NewContext(ctx, xl),
		config:   conf,
		handler:  handler,
		src_port: src,
		dst_port: dest,
	}

	t.src = helper.NewPingPong(ctx, src_config, src_stream)
	t.dst = helper.NewPingPong(ctx, dest_config, dest_stream)
	return t
}
