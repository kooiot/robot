package common

type Server interface {
	Run() error
	Init() error
}

type ServerHandler interface {
	Init(s Server) error
	AfterLogin(client *Client)
	BeforeLogout(client *Client)
}
