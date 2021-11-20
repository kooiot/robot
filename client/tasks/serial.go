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
	config   msg.SerialTask
	src_port *serial.SerialPort
	src      *helper.PingPong
	dst_port *serial.SerialPort
	dst      *helper.PingPong
}

func init() {
	RegisterTask("serial", NewSerialTask)
}

func (s *SerialTask) Start() error {
	return nil
}

func (s *SerialTask) Stop() error {
	return nil
}

func NewSerialTask(handler common.TaskHandler, option interface{}) common.Task {
	data, _ := json.Marshal(option)

	conf := msg.SerialTask{}
	json.Unmarshal(data, &conf)

	src_stream := port.NewStream()
	src, err := serial.NewSerial(src_stream, serial.Baudrate(conf.Baudrate))
	if err != nil {
		return nil
	}

	dest_stream := port.NewStream()
	dest, err := serial.NewSerial(dest_stream, serial.Baudrate(conf.Baudrate))
	if err != nil {
		return nil
	}

	return &SerialTask{
		config:   conf,
		src_port: src,
		src:      helper.NewPingPong(handler, helper.PingPongConfig{IsPing: true, Count: conf.Count, MaxMsgSize: conf.MaxMsgSize}, src_stream),
		dst_port: dest,
		dst:      helper.NewPingPong(handler, helper.PingPongConfig{IsPing: true, Count: conf.Count, MaxMsgSize: conf.MaxMsgSize}, dest_stream),
	}
}
