package server

import (
	"encoding/json"
	"strconv"

	"github.com/Allenxuxu/gev"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/net/protocol"
	"github.com/kooiot/robot/pkg/util/log"
	"github.com/kooiot/robot/server/config"
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
		req := msg.Login{}
		if err := json.Unmarshal(data, &req); err != nil {
			log.Error(err.Error())
		}
		log.Trace("received %s: %v", msgType, req)

		resp := msg.LoginResp{
			ClientID: req.ClientID,
			ID:       999,
			Reason:   "OK",
		}

		data, err := json.Marshal(resp)
		if err != nil {
			log.Error("failed encode resp: %s", err)
		} else {
			return protocol.PackMessage("login_resp", data)
		}
	case "logout":
		req := msg.Logout{}
		if err := json.Unmarshal(data, &req); err != nil {
			log.Error(err.Error())
		}
		log.Trace("received %s: %v", msgType, req)

		resp := &msg.Response{
			Content: "OK",
		}

		data, err := json.Marshal(resp)
		if err != nil {
			log.Error("failed encode resp: %s", err)
		} else {
			return protocol.PackMessage("logout_resp", data)
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
		gev.CustomProtocol(protocol.New()))
	if err != nil {
		panic(err)
	}
	handler.server = s

	return handler
}
