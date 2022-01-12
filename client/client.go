package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/kooiot/robot/client/config"
	"github.com/kooiot/robot/client/tasks"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/net/protocol"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type Client struct {
	ctx    context.Context
	cancel context.CancelFunc
	config *config.ClientConf
	conn   *Connection
	runner *tasks.Runner
}

func (c *Client) Run() error {
	xl := xlog.FromContextSafe(c.ctx)
	c.conn.OnOpen(func() {
		xl.Info("client opened")
		go c.OnRun()
	})

	c.conn.OnMessage(func(ctx interface{}, data []byte) (out interface{}) {
		return c.OnMessage(ctx, data)
	})

	return c.conn.Run()
}

func (c *Client) OnRun() {
	xl := xlog.FromContextSafe(c.ctx)
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
			xl.Info("Send login: %v", login)

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
			xl.Info("Send logout: %v", logout)

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

func (c *Client) SendResult(task *msg.Task, result *msg.TaskResult) error {
	xl := xlog.FromContextSafe(c.ctx)
	xl.Info("Send result: %#v", result)
	var buffer []byte

	str, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	buffer = protocol.PackMessage("task.result", str)

	_, err = c.conn.Write(buffer)
	return err
}

func (c *Client) SendTaskUpdate(task *msg.Task) error {
	xl := xlog.FromContextSafe(c.ctx)
	var buffer []byte
	xl.Info("Send task: %#v", task)

	str, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}
	buffer = protocol.PackMessage("task.update", str)

	_, err = c.conn.Write(buffer)
	return err
}

func (c *Client) OnMessage(ctx interface{}, data []byte) (out interface{}) {
	xl := xlog.FromContextSafe(c.ctx)
	msgType := ctx.(string)

	switch msgType {
	case "login_resp":
		msg := msg.LoginResp{}
		if err := json.Unmarshal(data, &msg); err != nil {
			xl.Info(err.Error())
		}
		xl.Info("%s: %v", msgType, msg)
	case "logout_resp":
		msg := msg.Response{}
		if err := json.Unmarshal(data, &msg); err != nil {
			xl.Info(err.Error())
		}
		xl.Info("%s: %v", msgType, msg)
	case "task":
		msg := msg.Task{}
		if err := json.Unmarshal(data, &msg); err != nil {
			xl.Info(err.Error())
		}
		xl.Info("%s: %v", msgType, msg)
		c.runner.Add(&msg, nil)
	default:
		xl.Info("unknown msg type %s", msgType)
	}

	return nil
}

func (c *Client) Close() {
	c.conn.Close()
	c.cancel()
}

func NewClient(cfg *config.ClientConf) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	cli := new(Client)
	cli.config = cfg
	cli.runner = tasks.NewRunner(ctx, &cfg.Runner, cli)
	cli.ctx = xlog.NewContext(ctx, xlog.New())
	cli.cancel = cancel

	addr := cfg.Common.Addr + ":" + strconv.Itoa(cfg.Common.Port)
	conn := NewConnection(ctx, addr)
	cli.conn = conn

	return cli
}
