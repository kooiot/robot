package client

import (
	"bufio"
	"context"
	"net"
	"time"

	"github.com/Allenxuxu/ringbuffer"
	"github.com/kooiot/robot/pkg/net/protocol"
)

type Connection struct {
	onOpenCallback    func()
	onMessageCallback func(ctx interface{}, data []byte) (out interface{})
	onErrorCallback   func(err error)

	ctx context.Context

	Conn      net.Conn
	Address   string
	Connected bool

	buffer *ringbuffer.RingBuffer
}

func (conn *Connection) OnOpen(f func()) {
	conn.onOpenCallback = f
}

func (conn *Connection) OnMessage(f func(ctx interface{}, data []byte) (out interface{})) {
	conn.onMessageCallback = f
}

func (conn *Connection) OnError(f func(err error)) {
	conn.onErrorCallback = f
}

func (conn *Connection) Close() error {
	return conn.Conn.Close()
}

func (conn *Connection) Write(message []byte) (n int, err error) {
	return conn.Conn.Write(message)
}

func (conn *Connection) WriteString(message string) (n int, err error) {
	return conn.Conn.Write([]byte(message))
}

func (conn *Connection) Run() error {
	c, err := net.DialTimeout("tcp", conn.Address, 3*time.Second)

	if err != nil {
		conn.onErrorCallback(err)
		return err
	} else {
		defer c.Close()
		conn.Conn = c

		conn.Connected = true
		conn.onOpenCallback()
		conn.read()
	}
	return nil
}

func (conn *Connection) read() {
	reader := bufio.NewReader(conn.Conn)

	for {
		buf := make([]byte, 1024)

		num, err := reader.Read(buf)

		if err != nil {
			conn.Close()
			conn.onErrorCallback(err)
			return
		}

		conn.buffer.Write(buf[:num])

		ctx, recvData := protocol.UnPacketMessage(conn.buffer)
		if ctx != nil || len(recvData) != 0 {
			sendData := conn.onMessageCallback(ctx, recvData)
			if sendData != nil {
				conn.Conn.Write(sendData.([]byte))
			}
		}
	}
}

func NewConnection(ctx context.Context, address string) *Connection {
	conn := &Connection{
		ctx:       ctx,
		Address:   address,
		Connected: false,
		buffer:    ringbuffer.New(0),
	}

	conn.OnOpen(func() {})
	conn.OnError(func(err error) {})
	conn.OnMessage(func(ctx interface{}, data []byte) (out interface{}) { return nil })

	return conn
}
