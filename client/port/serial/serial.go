package serial

import (
	"bufio"
	"errors"
	"sync"

	"github.com/Allenxuxu/toolkit/sync/atomic"
	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/pkg/util/log"
	"github.com/tarm/serial"
)

var ErrConnectionClosed = errors.New("serial closed")

type SerialPort struct {
	connected atomic.Bool
	handler   common.PortHandler
	serial    *serial.Port
	config    *serial.Config
	lock      sync.Mutex
}

func (s *SerialPort) Write(data []byte) error {
	if !s.connected.Get() {
		return ErrConnectionClosed
	}
	_, err := s.serial.Write(data)
	// TODO:
	return err
}

func (s *SerialPort) Open() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	log.Info("Serial open with :%#v", s.config)
	p, err := serial.OpenPort(s.config)
	if err != nil {
		log.Error("Serial open failed:%s", err.Error())
		s.handler.OnOpen(s, err)
		return err
	}
	s.serial = p
	s.handler.OnOpen(s, nil)
	s.connected.Set(true)

	go s.read()

	return nil
}

func (s *SerialPort) read() {
	reader := bufio.NewReader(s.serial)

	for {
		buf := make([]byte, 1024)
		s.lock.Lock()

		num, err := reader.Read(buf)

		if err != nil {
			if s.serial != nil {
				s.Close()
			}
			s.lock.Unlock()
			break
		}
		s.lock.Unlock()

		err = s.handler.OnMessage(buf[:num])
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func (s *SerialPort) Close() error {
	if !s.connected.Get() {
		return ErrConnectionClosed
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	p := s.serial
	s.serial = nil
	err := p.Close()
	s.connected.Set(false)
	s.handler.OnClose(err)
	return err
}

func NewSerial(handler common.PortHandler, opts ...Option) (*SerialPort, error) {
	if handler == nil {
		return nil, errors.New("handler is nil")
	}
	options := newOptions(opts...)
	port := SerialPort{handler: handler, lock: sync.Mutex{}}

	c := &serial.Config{
		Name:        options.Port,
		Baud:        options.Baudrate,
		Size:        options.DataBits,
		Parity:      serial.Parity(options.Parity),
		StopBits:    serial.StopBits(options.StopBits),
		ReadTimeout: options.ReadTimeout,
	}
	port.config = c
	return &port, nil
}
