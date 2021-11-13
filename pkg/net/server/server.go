package server

import (
	"flag"
	"log"
	"strconv"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/plugins/protobuf"
	pb "github.com/kooiot/robot/pkg/net/proto"
	"google.golang.org/protobuf/proto"
)

type server struct{}

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

func NewServer() {
	handler := new(server)
	var port int
	var loops int

	flag.IntVar(&port, "port", 1833, "server port")
	flag.IntVar(&loops, "loops", -1, "num loops")
	flag.Parse()

	s, err := gev.NewServer(handler,
		gev.Network("tcp"),
		gev.Address(":"+strconv.Itoa(port)),
		gev.NumLoops(loops),
		gev.CustomProtocol(&protobuf.Protocol{}))
	if err != nil {
		panic(err)
	}

	log.Println("server start")
	s.Start()
}
