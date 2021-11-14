package server

import (
	"log"
	"strconv"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/plugins/protobuf"
	pb "github.com/kooiot/robot/pkg/net/proto"
	"github.com/kooiot/robot/server/config"
	"google.golang.org/protobuf/proto"
)

type server struct {
	server *gev.Server
}

func (s *server) OnConnect(c *gev.Connection) {
	log.Println(" OnConnect ： ", c.PeerAddr())
}

func (s *server) OnMessage(c *gev.Connection, ctx interface{}, data []byte) (out interface{}) {
	msgType := ctx.(string)

	switch msgType {
	case "login":
		msg := &pb.Login{}
		if err := proto.Unmarshal(data, msg); err != nil {
			log.Println(err)
		}
		log.Println(msgType, msg)
	case "logout":
		msg := &pb.Logout{}
		if err := proto.Unmarshal(data, msg); err != nil {
			log.Println(err)
		}
		log.Println(msgType, msg)
	default:
		log.Println("unknown msg type", msgType)
	}

	return
}

func (s *server) OnClose(c *gev.Connection) {
	log.Println("OnClose")
}

func (s *server) Run() error {
	log.Println("server start")
	s.server.Start()
	return nil
}

func NewServer(cfg *config.ServerConf) *server {
	handler := new(server)

	s, err := gev.NewServer(handler,
		gev.Network("tcp"),
		gev.Address(cfg.Common.Bind+":"+strconv.Itoa(cfg.Common.Port)),
		gev.NumLoops(cfg.Common.Loops),
		gev.CustomProtocol(&protobuf.Protocol{}))
	if err != nil {
		panic(err)
	}
	handler.server = s

	return handler
}
