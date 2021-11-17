package tasks

import (
	"encoding/json"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/client/helper"
	"github.com/kooiot/robot/client/port"
	"github.com/kooiot/robot/client/port/serial"
)

type SerialTaskConfig struct {
	SrcPort    string `json:"src"`
	DestPort   string `json:"dst"`
	Baudrate   int    `json:"baudrate"`
	Count      int    `json:"count"`
	MaxMsgSize int    `json:"max_msg_size"`
}

type SerialTask struct {
	config   SerialTaskConfig
	src_port *serial.SerialPort
	src      *helper.PingPong
	dst_port *serial.SerialPort
	dst      *helper.PingPong
}

func (s *SerialTask) Start() error {
	return nil
}

func (s *SerialTask) Stop() error {
	return nil
}

func NewSerialTask(handler common.TaskHandler, config interface{}) *SerialTask {
	data, _ := json.Marshal(config)

	conf := SerialTaskConfig{}
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
