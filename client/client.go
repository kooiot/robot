package client

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"

	"github.com/Allenxuxu/gev/plugins/protobuf"
	"github.com/kooiot/robot/client/config"
	pb "github.com/kooiot/robot/pkg/net/proto"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	cfg *config.ClientConf
}

func (c *Client) Run() error {
	addr := c.cfg.Common.Addr + ":" + strconv.Itoa(c.cfg.Common.Port)
	conn, e := net.Dial("tcp", addr)
	if e != nil {
		log.Fatal(e)
	}
	defer conn.Close()

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

			data, err := proto.Marshal(msg)
			if err != nil {
				panic(err)
			}
			buffer = protobuf.PackMessage("login", data)
		case 1:
			msg := &pb.Logout{
				ClientId: name,
				Id:       "xxx",
			}

			data, err := proto.Marshal(msg)
			if err != nil {
				panic(err)
			}
			buffer = protobuf.PackMessage("logout", data)
		}

		_, err := conn.Write(buffer)
		if err != nil {
			panic(err)
		}
	}
}

func NewClient(cfg *config.ClientConf) *Client {
	cli := new(Client)
	cli.cfg = cfg

	return cli
}
