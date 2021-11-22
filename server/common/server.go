package common

type Server interface {
	Run() error
	Init() error
	ConfigDir() string
}

type ServerHandler interface {
	Init(s Server) error
	AfterLogin(client *Client)
	BeforeLogout(client *Client)
}
