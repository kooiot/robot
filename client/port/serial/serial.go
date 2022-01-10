package serial

import (
	"errors"
	"sync"

	"github.com/Allenxuxu/toolkit/sync/atomic"
	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/pkg/util/log"
	"go.bug.st/serial"
)

var ErrConnectionClosed = errors.New("serial closed")
var ErrAlreadyOpened = errors.New("serial already opened")

type SerialPort struct {
	connected atomic.Bool
	handler   common.PortHandler
	serial    serial.Port
	port      string
	mode      *serial.Mode
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
	if s.connected.Get() {
		return ErrAlreadyOpened
	}

	log.Info("Serial open %s with :%#v", s.port, s.mode)
	p, err := serial.Open(s.port, s.mode)
	if err != nil {
		log.Error("Serial open failed:%s", err.Error())
		s.handler.OnOpen(s, err)
		return err
	}

	s.lock.Lock()
	s.serial = p
	s.connected.Set(true)
	s.lock.Unlock()

	s.handler.OnOpen(s, nil)

	go s.read()

	return nil
}

func (s *SerialPort) read() {
	for {
		buf := make([]byte, 1024)
		if !s.connected.Get() {
			break
		}

		s.lock.Lock()
		port := s.serial
		s.lock.Unlock()
		if nil == port {
			break
		}

		num, err := port.Read(buf)

		if err != nil {
			if s.connected.Get() {
				log.Error("Serial closed! error: %s", err.Error())
				s.connected.Set(false)
				s.lock.Lock()
				s.serial = nil
				s.lock.Unlock()

				port.Close()
			}
			break
		}
		// log.Info("Serial read len: %d", num)

		err = s.handler.OnMessage(buf[:num])
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func (s *SerialPort) Close() error {
	log.Info("Serial close...")
	if !s.connected.Get() {
		return ErrConnectionClosed
	}

	s.connected.Set(false)

	s.lock.Lock()
	err := s.serial.Close()
	s.serial = nil
	s.lock.Unlock()

	s.handler.OnClose(nil)
	log.Info("Serial closed!!!")
	return err
}

func NewSerial(handler common.PortHandler, opts ...Option) (*SerialPort, error) {
	if handler == nil {
		return nil, errors.New("handler is nil")
	}
	options := newOptions(opts...)
	port := SerialPort{handler: handler, lock: sync.Mutex{}}
	port.port = options.Port

	if options.DataBits == 0 {
		options.DataBits = 8
	}

	c := &serial.Mode{
		BaudRate: options.Baudrate,
		DataBits: options.DataBits,
		Parity:   serial.Parity(options.Parity),
		StopBits: serial.StopBits(options.StopBits),
	}
	port.mode = c
	return &port, nil
}
