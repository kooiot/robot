package server

import (
	"strconv"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/plugins/protobuf"
	pb "github.com/kooiot/robot/pkg/net/proto"
	"github.com/kooiot/robot/pkg/util/log"
	"github.com/kooiot/robot/server/config"
	"google.golang.org/protobuf/proto"
)

type server struct {
	server *gev.Server
}

func (s *server) OnConnect(c *gev.Connection) {
	log.Info("client connected: %s", c.PeerAddr())
}

func (s *server) OnMessage(c *gev.Connection, ctx interface{}, data []byte) (out interface{}) {
	msgType := ctx.(string)

	switch msgType {
	case "login":
		msg := &pb.Login{}
		if err := proto.Unmarshal(data, msg); err != nil {
			log.Error(err.Error())
		}
		log.Trace("received %s: %v", msgType, msg)

		resp := &pb.LoginResp{
			ClientId: msg.ClientId,
			Id:       999,
			Reason:   "OK",
		}

		data, err := proto.Marshal(resp)
		if err != nil {
			log.Error("failed encode resp: %s", err)
		} else {
			return protobuf.PackMessage("login_resp", data)
		}
	case "logout":
		msg := &pb.Logout{}
		if err := proto.Unmarshal(data, msg); err != nil {
			log.Error(err.Error())
		}
		log.Trace("received %s: %v", msgType, msg)

		resp := &pb.Response{
			Content: "OK",
		}

		data, err := proto.Marshal(resp)
		if err != nil {
			log.Error("failed encode resp: %s", err)
		} else {
			return protobuf.PackMessage("logout_resp", data)
		}
	default:
		log.Error("unknown msg type %s", msgType)
	}

	return
}

func (s *server) OnClose(c *gev.Connection) {
	log.Info("client connection closed %s", c.PeerAddr())
}

func (s *server) Run() error {
	log.Info("server start")
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
