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
		length := int(buffer.PeekUint32())
		if buffer.Length() >= length+4 {
			buffer.Retrieve(4)

			typeLen := int(buffer.PeekUint16())
			buffer.Retrieve(2)

			typeByte := pbytes.GetLen(typeLen)
			_, _ = buffer.Read(typeByte)

			dataLen := length - 2 - typeLen
			data := make([]byte, dataLen)
			_, _ = buffer.Read(data)

			out = data
			ctx = string(typeByte)
			pbytes.Put(typeByte)
		}
	}
	return
}
