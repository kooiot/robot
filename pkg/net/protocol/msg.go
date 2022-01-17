package protocol

import (
	"encoding/binary"

	"github.com/Allenxuxu/ringbuffer"
	"github.com/gobwas/pool/pbytes"
)

// Message 数据帧定义
type Message struct {
	Len     uint32
	TypeLen uint16
	Type    string
	Data    []byte
}

// PackMessage 按自定义协议打包数据
func PackMessage(msgType string, data []byte) []byte {
	typeLen := len(msgType)
	length := len(data) + typeLen + 2

	ret := make([]byte, length+4)

	binary.BigEndian.PutUint32(ret, uint32(length))
	binary.BigEndian.PutUint16(ret[4:], uint16(typeLen))
	copy(ret[6:], msgType)
	copy(ret[6+typeLen:], data)

	return ret
}

// UnPacket ...
func UnPacketMessage(buffer *ringbuffer.RingBuffer) (ctx interface{}, out []byte) {
	if buffer.Length() > 6 {
		buf := pbytes.GetLen(4)
		defer pbytes.Put(buf)

		_, _ = buffer.VirtualRead(buf)
		length := binary.BigEndian.Uint32(buf)
		buffer.VirtualRevert()

		if buffer.Length() >= int(length)+4 {
			buffer.Retrieve(4)

			type_len := int(buffer.PeekUint16())
			buffer.Retrieve(2)

			type_bytes := pbytes.GetLen(type_len)
			defer pbytes.Put(type_bytes)

			_, _ = buffer.Read(type_bytes)

			data_len := int(length) - 2 - type_len
			data := make([]byte, data_len)
			_, _ = buffer.Read(data)

			out = data
			ctx = string(type_bytes)
		}
	}
	return
}
