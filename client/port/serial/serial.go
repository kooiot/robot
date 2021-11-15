package serial

import (
	"errors"

	"github.com/kooiot/robot/client/port"
	"github.com/kooiot/robot/pkg/util/log"
	"github.com/tarm/serial"
)

type SerialPort struct {
	handler port.Handler
	serial  *serial.Port
}

func (s *SerialPort) Write(data []byte) error {
	_, err := s.serial.Write(data)
	// TODO:
	return err
}

func NewPort(handler port.Handler, opts ...Option) (port.Port, error) {
	if handler == nil {
		return nil, errors.New("handler is nil")
	}
	options := newOptions(opts...)
	port := SerialPort{handler: handler}

	c := &serial.Config{
		Name:        options.Port,
		Baud:        options.Baudrate,
		Size:        options.DataBits,
		Parity:      serial.Parity(options.Parity),
		StopBits:    serial.StopBits(options.StopBits),
		ReadTimeout: options.ReadTimeout,
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Info(err.Error())
		return nil, err
	}
	port.serial = s
	return &port, nil
}
