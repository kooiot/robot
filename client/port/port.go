package port

type Port interface {
	Write(c []byte) error
}

// Port Handler 注册接口
type Handler interface {
	OnConnect(port *Port)
	OnOpen(error)
	OnClose(error)
	OnMessage([]byte) error
}
