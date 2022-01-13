package client

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/kooiot/robot/client/config"
	"github.com/kooiot/robot/client/tasks"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/net/protocol"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type Client struct {
	ctx            context.Context
	cancel         context.CancelFunc
	lock           sync.Mutex
	config         *config.ClientConf
	conn           *Connection
	runner         *tasks.Runner
	client_id      int32
	last_heartbeat time.Time
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

	c.lock.Lock()
	c.client_id = -1
	c.last_heartbeat = time.Now().UTC()
	c.lock.Unlock()

	timeout := 100 * time.Millisecond
	for {
		time.Sleep(timeout)

		if c.client_id < 0 {
			err := c.send_login()
			if err != nil {
				xl.Error("Failed to login: %s", err.Error())
				timeout = timeout * 2
				if timeout > time.Second*30 {
					timeout = 100 * time.Millisecond
				}
			} else {
				timeout = time.Second
				time.Sleep(3 * time.Second)
			}
		} else {
			if time.Since(c.last_heartbeat) > 60*time.Second {
				err := c.send_heartbeat()
				if err == nil {
					c.last_heartbeat = c.last_heartbeat.Add(3 * time.Second)
				}
			}
		}
	}
}

func (c *Client) send_message(msg_type string, msg_data interface{}) error {
	xl := xlog.FromContextSafe(c.ctx)
	xl.Debug("Send %s: %#v", msg_type, msg_data)

	data, err := json.Marshal(msg_data)
	if err != nil {
		panic(err)
	}
	buffer := protocol.PackMessage(msg_type, data)

	c.lock.Lock()
	defer c.lock.Unlock()
	_, err = c.conn.Write(buffer)

	// xl.Info("Send %s done", msg_type)

	return err
}

func (c *Client) send_login() error {
	req := &msg.Login{
		ClientID: c.config.Common.ClientID,
		User:     c.config.Common.User,
		Passwd:   c.config.Common.Password,
		Hostname: "Hostname",
		Hardware: "ARM v7",
		System:   "OpenWRT",
	}
	return c.send_message("login", &req)
}

func (c *Client) send_logout() error {
	req := &msg.Logout{
		ClientID: c.config.Common.ClientID,
		ID:       c.client_id,
	}
	return c.send_message("logout", &req)
}

func (c *Client) send_heartbeat() error {
	req := &msg.HeartBeat{
		ID:   c.client_id,
		Time: time.Now().UTC().Unix(),
	}
	return c.send_message("heartbeat", &req)
}

func (c *Client) SendResult(task *msg.Task, result *msg.TaskResult) error {
	return c.send_message("task.result", result)
}

func (c *Client) SendTaskUpdate(task *msg.Task) error {
	return c.send_message("task.result", task)
}

func (c *Client) OnMessage(ctx interface{}, data []byte) (out interface{}) {
	xl := xlog.FromContextSafe(c.ctx)
	msgType := ctx.(string)
	xl.Debug("On Message: %s", msgType)

	c.lock.Lock()
	defer c.lock.Unlock()

	switch msgType {
	case "login_resp":
		req := msg.LoginResp{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Info(err.Error())
		}
		xl.Debug("%s: %v", msgType, req)
		c.client_id = req.ID
	case "logout_resp":
		req := msg.Response{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Info(err.Error())
		}
		xl.Debug("%s: %v", msgType, req)
	case "heartbeat":
		req := msg.HeartBeat{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Info(err.Error())
		}
		c.last_heartbeat = time.Now().UTC()
		xl.Debug("%s: %v", msgType, req)
	case "task":
		req := msg.Task{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Info(err.Error())
		}
		xl.Info("%s: %v", msgType, req)
		go c.runner.Add(&req, nil)
	default:
		xl.Info("unknown msg type %s", msgType)
	}

	return nil
}

func (c *Client) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.send_logout()
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
