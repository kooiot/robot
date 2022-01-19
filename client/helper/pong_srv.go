package helper

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"time"

	"github.com/Allenxuxu/ringbuffer"
	"github.com/Allenxuxu/toolkit/sync/atomic"
	"github.com/kooiot/robot/client/port"
	"github.com/kooiot/robot/pkg/util/xlog"
	"github.com/npat-efault/crc16"
)

type PongSrv struct {
	ctx    context.Context
	stream *port.Stream
	stop   atomic.Bool
}

func (c *PongSrv) Start() {
	c.stop.Set(false)
	go c.run()
}

func (c *PongSrv) Stop() {
	c.stop.Set(true)
}

func (c *PongSrv) IsRunning() bool {
	return !c.stop.Get()
}

func (c *PongSrv) run() {
	xl := xlog.FromContextSafe(c.ctx)

	// Stop stream
	defer c.stream.Stop()
	defer c.stop.Set(true)

	xl.Debug("PongSrv test start.")

	reqMsg := make([]byte, 0)

	for {
		if c.stop.Get() {
			break
		}
		recv, err := c.stream.Request(reqMsg, func(buffer *ringbuffer.RingBuffer) ([]byte, error) {
			// Try to find the SK
			for {
				data := buffer.Bytes()
				if buffer.Length() < hdr_len+end_len {
					return nil, nil
				}

				// Find SK
				sk_i := bytes.Index(data, SK)
				if sk_i > 0 {
					xl.Error("SK checking error %d", sk_i)
					buffer.Retrieve(sk_i)
				} else if sk_i < 0 {
					xl.Error("SK Find None %d", len(data))
					buffer.Retrieve(len(data) - 3)
					return nil, nil
				} else {
					break
				}
			}
			data := buffer.Bytes()

			// Read len
			data_len := binary.BigEndian.Uint16(data[3:5])
			// xl.Info("PongSrv recv msg len: %d", data_len)

			if len(data) < hdr_len+end_len+int(data_len) {
				// xl.Info("len:%d data_len:%d", len(data), hdr_len+end_len+int(data_len))
				return nil, nil
			}

			if !bytes.Equal(EK, data[hdr_len+int(data_len)+2:hdr_len+int(data_len)+5]) {
				xl.Error("EK Check Error %d", len(data))
				buffer.Retrieve(1)
				return nil, nil
			} else {
				// Retrieve buffer
				buffer.Retrieve(hdr_len + end_len + int(data_len))
				// xl.Info("left size: %d", buffer.Length())

				h := crc16.New(crc16.Modbus)
				h.Write(data[hdr_len : hdr_len+int(data_len)])
				crc_16 := binary.BigEndian.Uint16(data[hdr_len+int(data_len) : hdr_len+int(data_len)+2])
				if crc_16 != h.Sum16() {
					xl.Error("crc checking error")
					return nil, errors.New("crc checking error")
				} else {
					// xl.Info("crc checking done: %x", crc_16)
				}
				return data[0 : hdr_len+end_len+int(data_len)], nil
			}
		}, time.Millisecond*1000)

		if err != nil {
			xl.Error("resp error: %s", err.Error())
		}

		// xl.Info("write back len: %d", len(recv))
		reqMsg = recv
	}

	xl.Debug("PongSrv test finished")
}

func NewPongSrv(ctx context.Context, stream *port.Stream) *PongSrv {
	xl := xlog.FromContextSafe(ctx).Spawn().AppendPrefix("PongSrv")
	o := PongSrv{
		ctx:    xlog.NewContext(ctx, xl),
		stream: stream,
		stop:   atomic.Bool{},
	}

	return &o
}
