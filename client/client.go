package client

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/kooiot/robot/client/config"
	pb "github.com/kooiot/robot/pkg/net/proto"
	"github.com/kooiot/robot/pkg/util/log"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	cfg   *config.ClientConf
	conn  *Connection
	proto *Protocol
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
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		name := text[:len(text)-1]

		switch rand.Int() % 2 {
		case 0:
			msg := &pb.Login{
				ClientId: name,
				User:     "User",
				Passwd:   "Passwd",
				Hostname: "Hostname",
				Hardware: "ARM v7",
				System:   "OpenWRT",
			}
			log.Info("Send login: %v", msg)

			data, err := proto.Marshal(msg)
			if err != nil {
				panic(err)
			}
			buffer = c.proto.Packet("login", data)
		case 1:
			msg := &pb.Logout{
				ClientId: name,
				Id:       "xxx",
			}
			log.Info("Send logout: %v", msg)

			data, err := proto.Marshal(msg)
			if err != nil {
				panic(err)
			}
			buffer = c.proto.Packet("logout", data)
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
		msg := &pb.LoginResp{}
		if err := proto.Unmarshal(data, msg); err != nil {
			log.Info(err.Error())
		}
		log.Info("%s: %v", msgType, msg)
	case "logout_resp":
		msg := &pb.Response{}
		if err := proto.Unmarshal(data, msg); err != nil {
			log.Info(err.Error())
		}
		log.Info("%s: %v", msgType, msg)
	default:
		log.Info("unknown msg type %s", msgType)
	}

	return nil
}

func NewClient(cfg *config.ClientConf) *Client {
	cli := new(Client)
	cli.cfg = cfg

	cli.proto = NewProtocol()
	addr := cfg.Common.Addr + ":" + strconv.Itoa(cfg.Common.Port)
	conn := NewConnection(addr, cli.proto)
	cli.conn = conn

	return cli
}
