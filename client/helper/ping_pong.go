package helper

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/rand"
	"time"

	"github.com/Allenxuxu/ringbuffer"
	"github.com/Allenxuxu/toolkit/sync/atomic"
	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/client/port"
	"github.com/kooiot/robot/pkg/util/log"
	"github.com/npat-efault/crc16"
)

type PingPongConfig struct {
	IsPing     bool `json:"is_ping"`
	Count      int  `json:"count"`
	MaxMsgSize int  `json:"max_msg_size"`
}

type PingPong struct {
	Config  PingPongConfig      `json:"config"`
	Result  common.StreamResult `json:"result"`
	task    common.Task
	handler common.TaskHandler
	stream  *port.Stream
	stop    atomic.Bool
}

var SK = []byte("AAA")
var EK = []byte("FFF")
var hdr_len = 3 + 2
var end_len = 3 + 2

func (c *PingPong) Start() {
	c.stop.Set(false)
	go c.run()
}

func (c *PingPong) Stop() {
	c.stop.Set(true)
}

func (c *PingPong) genMsg() []byte {
	var b bytes.Buffer // A Buffer needs no initialization.

	// Write packet header
	b.Write(SK)

	data_len := rand.Intn(c.Config.MaxMsgSize)
	if data_len < 16 {
		data_len = 16 + rand.Intn(16)
	}
	// log.Info("PingPong generate msg len: %d", data_len)

	buf := make([]byte, data_len)
	for i := 0; i < data_len; i++ {
		b := rand.Intn(256)
		buf[i] = byte(b)
	}
	// Write Random string length and content
	binary.Write(&b, binary.BigEndian, uint16(data_len))
	b.Write(buf)

	h := crc16.New(crc16.Modbus)
	h.Write(buf)
	// Append CRC16
	binary.Write(&b, binary.BigEndian, h.Sum16())

	b.Write(EK)

	if b.Len() != data_len+hdr_len+end_len {
		panic("Message len incorrect")
	}

	return b.Bytes()
}

func (c *PingPong) run() error {
	c.Result.SendSpeed = 0
	c.Result.RecvSpeed = 0
	if c.Config.IsPing {
		time.Sleep(time.Millisecond * 200)
	}
	send_total := 0
	recv_total := 0
	err_total := 0
	begin_time := time.Now()

	log.Info("PingPong test start. config: %#v", c.Config)

	msg := make([]byte, 0)

	for i := 0; i < c.Config.Count; i++ {
		if c.stop.Get() {
			break
		}
		if c.Config.IsPing {
			msg = c.genMsg()
		}
		send_total += len(msg)
		recv, err := c.stream.Request(msg, func(buffer *ringbuffer.RingBuffer) ([]byte, error) {
			// Try to find the SK
			for {
				data := buffer.Bytes()
				if buffer.Length() < hdr_len+end_len {
					return nil, nil
				}

				// Find SK
				sk_i := bytes.Index(data, SK)
				if sk_i > 0 {
					log.Error("SK checking error %d", sk_i)
					err_total += sk_i
					buffer.Retrieve(sk_i)
				} else if sk_i < 0 {
					log.Error("SK Find None %d", len(data))
					err_total += len(data) - 3
					buffer.Retrieve(len(data) - 3)
					return nil, nil
				} else {
					break
				}
			}
			data := buffer.Bytes()

			// Read len
			data_len := binary.BigEndian.Uint16(data[3:5])
			// log.Info("PingPong recv msg len: %d", data_len)

			if len(data) < hdr_len+end_len+int(data_len) {
				// log.Info("len:%d data_len:%d", len(data), hdr_len+end_len+int(data_len))
				return nil, nil
			}

			if !bytes.Equal(EK, data[hdr_len+int(data_len)+2:hdr_len+int(data_len)+5]) {
				log.Error("EK Check Error %d", len(data))
				err_total += 1
				buffer.Retrieve(1)
				return nil, nil
			} else {
				// Retrieve buffer
				buffer.Retrieve(hdr_len + end_len + int(data_len))
				// log.Info("left size: %d", buffer.Length())

				h := crc16.New(crc16.Modbus)
				h.Write(data[hdr_len : hdr_len+int(data_len)])
				crc_16 := binary.BigEndian.Uint16(data[hdr_len+int(data_len) : hdr_len+int(data_len)+2])
				if crc_16 != h.Sum16() {
					err_total += hdr_len + end_len + int(data_len)
					log.Error("crc checking error")
					return nil, errors.New("crc checking error")
				} else {
					// log.Info("crc checking done: %x", crc_16)
				}
				return data[0 : hdr_len+end_len+int(data_len)], nil
			}
		}, time.Millisecond*1000)

		if err != nil {
			log.Error("resp error:%s", err.Error())
			c.Result.Failed += 1
		} else {
			c.Result.Passed += 1
			recv_total += len(recv)
		}

		c.Result.Count += 1
		if !c.Config.IsPing {
			// log.Info("write back len: %d", len(recv))
			msg = recv
		}
	}
	if len(msg) > 0 && !c.Config.IsPing {
		// Write last recv
		c.stream.Send(msg, 100)
	}

	c.Result.ErrBytes = err_total
	c.Result.RecvBytes = recv_total + err_total
	c.Result.SendBytes = send_total
	c.Result.SendSpeed = float64(send_total) / time.Since(begin_time).Seconds()
	c.Result.RecvSpeed = float64(recv_total) / time.Since(begin_time).Seconds()

	log.Info("PingPong test finished: %#v", c.Result)
	result := common.TaskResult{
		Result: true,
		Error:  "Done",
		Info:   c.Result,
	}
	c.handler.OnResult(c.task, result)

	// Stop stream
	defer c.stream.Stop()

	return nil
}

func NewPingPong(task common.Task, handler common.TaskHandler, c PingPongConfig, stream *port.Stream) *PingPong {
	if c.MaxMsgSize == 0 {
		c.MaxMsgSize = 512
	}

	o := PingPong{
		Config:  c,
		task:    task,
		handler: handler,
		stream:  stream,
		stop:    atomic.Bool{},
	}

	return &o
}
