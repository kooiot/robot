package port

import (
	"errors"
	"sync"
	"time"

	"github.com/Allenxuxu/ringbuffer"
	"github.com/kooiot/robot/client/common"
)

type StreamParser func(*ringbuffer.RingBuffer) ([]byte, error)

type Stream struct {
	port   common.Port
	lock   sync.Mutex
	buffer *ringbuffer.RingBuffer
	// channel chan []byte
	// parser StreamParser
}

func (s *Stream) OnOpen(port common.Port, err error) {
	if err != nil {
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	s.buffer.RetrieveAll()
	s.port = port
}

func (s *Stream) OnClose(error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.port = nil
}

func (s *Stream) OnMessage(data []byte) error {
	left := len(data)
	for {
		s.lock.Lock()
		n, err := s.buffer.Write(data)
		s.lock.Unlock()

		if err != nil {
			return err
		}
		left -= n
		if left == 0 {
			break
		} else {
			time.Sleep(5 * time.Millisecond)
		}
	}
	return nil
}

func (s *Stream) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.port != nil {
		s.port.Close()
	}
}

func (s *Stream) isOpen() bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.port != nil
}

func (s *Stream) Send(data []byte, timeout time.Duration) error {
	begin := time.Now()
	for {
		if s.isOpen() {
			break
		}
		time.Sleep(10 * time.Millisecond)
		if time.Since(begin) > timeout {
			return errors.New("timeout")
		}
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	return s.port.Write(data)
}

func (s *Stream) Request(data []byte, parser StreamParser, timeout time.Duration) ([]byte, error) {
	err := s.Send(data, timeout)

	if err != nil {
		return nil, err
	}

	begin := time.Now()

	var msg []byte
	last_len := 0
	for {
		{
			s.lock.Lock()
			if s.buffer.Length() > last_len {
				msg, err = parser(s.buffer)
				if err != nil {
					s.lock.Unlock()
					break
				}
				if len(msg) > 0 {
					s.lock.Unlock()
					break
				}
				last_len = s.buffer.Length()
			}
			s.lock.Unlock()
		}
		time.Sleep(10 * time.Millisecond)
		if time.Since(begin) > timeout {
			if last_len > 0 {
				s.lock.Lock()
				s.buffer.RetrieveAll()
				s.lock.Unlock()
			}
			return nil, errors.New("timeout")
		}
	}

	return msg, nil
}

func NewStream() *Stream {
	s := &Stream{
		lock:   sync.Mutex{},
		buffer: ringbuffer.New(0),
	}
	// s.channel = make(chan []byte)
	return s
}
