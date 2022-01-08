package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/kooiot/robot/client/config"
	"github.com/kooiot/robot/client/tasks"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/net/protocol"
	"github.com/kooiot/robot/pkg/util/log"
)

type Client struct {
	config *config.ClientConf
	conn   *Connection
	runner *tasks.Runner
}

func (c *Client) Run() error {
	c.conn.OnOpen(func() {
		log.Info("client opened")
		go c.OnRun()
	})

	c.conn.OnMessage(func(ctx interface{}, data []byte) (out interface{}) {
		return c.OnMessage(ctx, data)
	})

	return c.conn.Run()
}

func (c *Client) OnRun() {
	var buffer []byte
	login := false
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		name := text[:len(text)-1]

		if !login {
			login := &msg.Login{
				ClientID: c.config.Common.ClientID,
				User:     c.config.Common.User,
				Passwd:   c.config.Common.Password,
				Hostname: "Hostname",
				Hardware: "ARM v7",
				System:   "OpenWRT",
			}
			log.Info("Send login: %v", login)

			data, err := json.Marshal(login)
			if err != nil {
				panic(err)
			}
			buffer = protocol.PackMessage("login", data)
		} else {
			logout := &msg.Logout{
				ClientID: name,
				ID:       0,
			}
			log.Info("Send logout: %v", logout)

			data, err := json.Marshal(logout)
			if err != nil {
				panic(err)
			}
			buffer = protocol.PackMessage("login", data)
		}

		_, err := c.conn.Write(buffer)
		if err != nil {
			panic(err)
		}
	}
}

func (c *Client) OnMessage(ctx interface{}, data []byte) (out interface{}) {
	msgType := ctx.(string)

	switch msgType {
	case "login_resp":
		msg := msg.LoginResp{}
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Info(err.Error())
		}
		log.Info("%s: %v", msgType, msg)
	case "logout_resp":
		msg := msg.Response{}
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Info(err.Error())
		}
		log.Info("%s: %v", msgType, msg)
	case "task":
		msg := msg.Task{}
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Info(err.Error())
		}
		log.Info("%s: %v", msgType, msg)
		c.runner.Add(&msg, nil)
	default:
		log.Info("unknown msg type %s", msgType)
	}

	return nil
}

func NewClient(cfg *config.ClientConf) *Client {
	cli := new(Client)
	cli.config = cfg
	cli.runner = tasks.NewRunner()

	addr := cfg.Common.Addr + ":" + strconv.Itoa(cfg.Common.Port)
	conn := NewConnection(addr)
	cli.conn = conn

	return cli
}
