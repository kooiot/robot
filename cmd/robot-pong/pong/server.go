package pong

import (
	"context"

	"github.com/kooiot/robot/client/helper"
	"github.com/kooiot/robot/client/port"
	"github.com/kooiot/robot/client/port/serial"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type SerialConfig struct {
	Port     string `mapstructure:"port" json:"port"`
	Baudrate int    `mapstructure:"baudrate" json:"baudrate"`
}

type PongServer struct {
	ctx    context.Context
	config SerialConfig

	serial *serial.SerialPort
	srv    *helper.PongSrv
}

func (s *PongServer) Run() {
	xl := xlog.FromContextSafe(s.ctx)
	err := s.serial.Open()
	if err != nil {
		xl.Error(err.Error())
	}
	s.srv.Start()
}

func (s *PongServer) IsRunning() bool {
	return s.srv.IsRunning()
}

func (s *PongServer) Abort() error {
	s.srv.Stop()
	return nil
}

func NewPongServer(ctx context.Context, cfg SerialConfig) *PongServer {
	stream := port.NewStream()
	ser, err := serial.NewSerial(stream, serial.Port(cfg.Port), serial.Baudrate(cfg.Baudrate))
	if err != nil {
		return nil
	}
	srv := helper.NewPongSrv(ctx, stream)

	t := &PongServer{
		ctx:    ctx,
		config: cfg,
		serial: ser,
		srv:    srv,
	}

	return t
}
