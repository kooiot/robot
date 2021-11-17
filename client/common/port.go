package common

type Port interface {
	Open() error
	Close() error
	Write([]byte) error
}

// Port Handler 注册接口
type PortHandler interface {
	OnOpen(Port, error)
	OnClose(error)
	OnMessage([]byte) error
}
