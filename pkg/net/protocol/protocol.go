//go:build !windows
// +build !windows

package protocol

import (
	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/ringbuffer"
)

// Protocol protobuf
type Protocol struct {
}

// New 创建 protobuf Protocol
func New() *Protocol {
	return &Protocol{}
}

// UnPacket ...
func (p *Protocol) UnPacket(c *gev.Connection, buffer *ringbuffer.RingBuffer) (ctx interface{}, out []byte) {
	return UnPacketMessage(buffer)
}

// Packet ...
func (p *Protocol) Packet(c *gev.Connection, data interface{}) []byte {
	return data.([]byte)
}
