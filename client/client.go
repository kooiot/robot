package client

import (
	"context"
	"encoding/json"
	"net"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/kooiot/robot/client/config"
	"github.com/kooiot/robot/client/tasks"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/net/protocol"
	"github.com/kooiot/robot/pkg/util/shutdown"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type Client struct {
	ctx context.Context

	xl *xlog.Logger

	config *config.ClientConf

	// Client connection
	conn *Connection

	// write to this channel to write the message sent to server
	send_chn chan (*msg.Message)
	// read from this channel to get the next message sent by server
	read_chn chan (*msg.Message)
	// goroutines can block by reading from this channel, it will be closed only in reader() when control connection is closed
	closed_chn chan struct{}
	// closing done event
	closed_done_chn chan struct{}
	// connection ready event
	connected_chn chan struct{}

	// Task Runner
	runner *tasks.Runner
	// Client ID
	client_id int32
	// Heartbeat
	last_heartbeat time.Time

	readerShutdown     *shutdown.Shutdown
	writerShutdown     *shutdown.Shutdown
	msgHandlerShutdown *shutdown.Shutdown
}

func (c *Client) newConn() (*Connection, error) {
	xl := c.xl

	cfg := c.config
	addr := net.JoinHostPort(cfg.Common.Addr, strconv.Itoa(cfg.Common.Port))
	xl.Info("Connect to %s", addr)

	conn := NewConnection(c.ctx, addr)

	conn.OnOpen(func() {
		err := c.send_login()
		if err != nil {
			xl.Error("Failed to login: %s", err.Error())
			c.conn.Close()
		}
	})

	conn.OnMessage(func(ctx interface{}, data []byte) (out interface{}) {
		xl.Debug("conn.OnMessage %v", ctx)
		c.read_chn <- &msg.Message{
			CTX:  ctx,
			Data: data,
		}
		return nil
	})
	conn.OnError(func(err error) {
		xl.Error("Connection closed: %s", err.Error())
		close(c.closed_chn)
	})

	return conn, nil
}

func (c *Client) Run() error {
	conn, err := c.newConn()
	if err != nil {
		return err
	}

	c.conn = conn

	go c.worker()

	return nil
}

// reader read all messages from frps and send to readCh
func (c *Client) reader() {
	xl := c.xl
	defer func() {
		if err := recover(); err != nil {
			xl.Error("panic error: %v", err)
			xl.Error(string(debug.Stack()))
		}
	}()
	defer c.readerShutdown.Done()
	defer close(c.closed_chn)

	xl.Info("client connection run start")
	c.conn.Run()
	xl.Error("client connection run finished!")
}

// writer writes messages got from sendCh to frps
func (c *Client) writer() {
	xl := c.xl
	defer c.writerShutdown.Done()
	for {
		m, ok := <-c.send_chn
		if !ok {
			xl.Info("send channel closed!")
			break
		} else {
			if err := c.write_conn(m); err != nil {
				xl.Error("Send message failed: %s", err.Error())
			}
		}
	}
}

func (c *Client) worker() {
	go c.msgHandler()
	go c.reader()
	go c.writer()

	<-c.closed_chn

	close(c.read_chn)
	c.readerShutdown.WaitDone()
	c.msgHandlerShutdown.WaitDone()

	close(c.send_chn)
	c.writerShutdown.WaitDone()

	close(c.closed_done_chn)
}

func (c *Client) msgHandler() {
	xl := c.xl
	defer func() {
		if err := recover(); err != nil {
			xl.Error("panic error: %v", err)
			xl.Error(string(debug.Stack()))
		}
	}()
	defer c.msgHandlerShutdown.Done()

	hbSend := time.NewTicker(60 * time.Second)
	defer hbSend.Stop()
	hbCheck := time.NewTicker(time.Second)
	defer hbCheck.Stop()

	c.last_heartbeat = time.Now()

	for {
		select {
		case <-hbSend.C:
			c.send_heartbeat()
		case <-hbCheck.C:
			if time.Since(c.last_heartbeat) > 90*time.Second {
				xl.Warn("heartbeat timeout")
				// let reader() stop
				c.conn.Close()
				return
			}
		case m, ok := <-c.read_chn:
			if !ok {
				return
			}
			c.OnMessage(m.CTX, m.Data)
		}
	}
}

func (c *Client) write_conn(m *msg.Message) error {
	buffer := protocol.PackMessage(m.CTX.(string), m.Data)

	_, err := c.conn.Write(buffer)

	return err
}

func (c *Client) send_message(msg_type string, msg_data interface{}) error {
	xl := c.xl
	xl.Debug("Send %s: %#v", msg_type, msg_data)

	data, err := json.Marshal(msg_data)
	if err != nil {
		xl.Error("JSON.Marshal failure: %s", err.Error())
		return err
	}
	c.send_chn <- &msg.Message{CTX: msg_type, Data: data}
	return nil
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

func (c *Client) SendResult(result *msg.TaskResult) error {
	return c.send_message("task.result", result)
}

func (c *Client) SendTaskUpdate(task *msg.Task) error {
	return c.send_message("task.update", task)
}

func (c *Client) PrintDone() {
	xl := c.xl

	for i := 0; i < 3; i++ {
		xl.Debug("ClientID: %s\r\n", c.config.Common.ClientID)
	}
}

func (c *Client) OnMessage(ctx interface{}, data []byte) (out interface{}) {
	xl := c.xl
	msgType := ctx.(string)
	xl.Debug("On Message: %s", msgType)

	switch msgType {
	case "login.resp":
		req := msg.LoginResp{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Error("JSON.Unmarshal error: %s", err.Error())
		}
		xl.Debug("%s: %#v", msgType, req)
		c.client_id = req.ID
		close(c.connected_chn)
	case "logout.resp":
		req := msg.Response{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Error("JSON.Unmarshal error: %s", err.Error())
		}
		xl.Debug("%s: %#v", msgType, req)
	case "heartbeat":
		req := msg.HeartBeat{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Error("JSON.Unmarshal error: %s", err.Error())
		}
		c.last_heartbeat = time.Now()
		xl.Debug("%s: %#v", msgType, req)
	case "task":
		req := msg.Task{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Error("JSON.Unmarshal error: %s", err.Error())
		}
		xl.Debug("%s: %#v", msgType, req)
		c.runner.Add(req, nil)
	default:
		xl.Error("unknown msg type %s", msgType)
	}

	return nil
}

func (c *Client) ConnectedChn() <-chan struct{} {
	return c.connected_chn
}

func (c *Client) ClosedDoneChn() <-chan struct{} {
	return c.closed_done_chn
}

func (c *Client) Close() {
	c.send_logout()
	time.Sleep(time.Second)
	c.conn.Close()
}

func NewClient(cfg *config.ClientConf, ctx context.Context) *Client {
	cli := &Client{
		config:             cfg,
		ctx:                ctx,
		xl:                 xlog.FromContextSafe(ctx),
		send_chn:           make(chan *msg.Message, 100),
		read_chn:           make(chan *msg.Message, 100),
		closed_chn:         make(chan struct{}),
		closed_done_chn:    make(chan struct{}),
		connected_chn:      make(chan struct{}),
		readerShutdown:     shutdown.New(),
		writerShutdown:     shutdown.New(),
		msgHandlerShutdown: shutdown.New(),
	}
	cli.runner = tasks.NewRunner(ctx, cli)

	return cli
}
