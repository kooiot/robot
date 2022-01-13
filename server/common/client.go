package common

import (
	"time"

	"github.com/Allenxuxu/gev"
	"github.com/kooiot/robot/pkg/net/msg"
)

type Client struct {
	ID            int32
	Info          *msg.Login
	Conn          *gev.Connection
	LastHeartbeat time.Time
}

func NewClient(id int32, conn *gev.Connection, login *msg.Login) *Client {
	return &Client{
		ID:            id,
		Info:          login,
		Conn:          conn,
		LastHeartbeat: time.Now(),
	}
}
