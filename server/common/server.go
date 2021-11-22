package common

import "github.com/Allenxuxu/gev"

type Server interface {
	Run() error
	Init() error
	ConfigDir() string
}

type ServerHandler interface {
	Init(s Server) error
	AfterLogin(conn *gev.Connection, client *Client)
	BeforeLogout(conn *gev.Connection, client *Client)
}
