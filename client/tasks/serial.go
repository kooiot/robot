package tasks

import (
	"encoding/json"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/client/helper"
	"github.com/kooiot/robot/client/port"
	"github.com/kooiot/robot/client/port/serial"
	"github.com/kooiot/robot/pkg/net/msg"
)

type SerialTask struct {
	info    *msg.Task
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

func (s *SerialTask) Start() error {
	err := s.src_port.Open()
	if err != nil {
		return err
	}
	err = s.dst_port.Open()
	if err != nil {
		return err
	}
	s.src.Start()
	s.dst.Start()
	return nil
}

func (s *SerialTask) Stop() error {
	return nil
}

func NewSerialTask(handler common.TaskHandler, info *msg.Task, parent common.Task) common.Task {
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
	dest_config := helper.PingPongConfig{IsPing: true, Count: conf.Count, MaxMsgSize: conf.MaxMsgSize}

	t := &SerialTask{
		info:     info,
		config:   conf,
		handler:  handler,
		src_port: src,
		dst_port: dest,
	}

	t.src = helper.NewPingPong(t, handler, src_config, src_stream)
	t.dst = helper.NewPingPong(t, handler, dest_config, dest_stream)
	return t
}
