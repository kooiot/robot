package server

import (
	"encoding/json"
	"path"
	"strconv"
	"time"

	"github.com/Allenxuxu/gev"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/net/protocol"
	"github.com/kooiot/robot/pkg/util/log"
	"github.com/kooiot/robot/server/common"
	"github.com/kooiot/robot/server/config"
	"github.com/kooiot/robot/server/tasks"
)

type Server struct {
	config      *config.ServerConf
	config_file string
	server      *gev.Server
	handlers    []common.ServerHandler
	clients     []common.Client
}

func (s *Server) OnConnect(c *gev.Connection) {
	log.Info("client connected: %s", c.PeerAddr())
}

func (s *Server) AfterLogin(conn *gev.Connection, client *common.Client) {
	time.Sleep(1 * time.Second)

	log.Info("AfterLogin %s", conn.PeerAddr())
	for _, h := range s.handlers {
		h.AfterLogin(conn, client)
	}
}

func (s *Server) HandleLogin(c *gev.Connection, req *msg.Login) interface{} {
	resp := msg.LoginResp{
		ClientID: req.ClientID,
		ID:       999,
		Reason:   "OK",
	}
	client := common.Client{
		Info: *req,
	}
	s.clients = append(s.clients, client)

	data, err := json.Marshal(resp)
	if err != nil {
		log.Error("failed encode resp: %s", err)
		return nil
	} else {
		go s.AfterLogin(c, &client)
		return protocol.PackMessage("login_resp", data)
	}
}

func (s *Server) OnMessage(c *gev.Connection, ctx interface{}, data []byte) interface{} {
	msgType := ctx.(string)

	switch msgType {
	case "login":
		req := msg.Login{}
		if err := json.Unmarshal(data, &req); err != nil {
			log.Error(err.Error())
		}
		log.Trace("received %s: %v", msgType, req)
		return s.HandleLogin(c, &req)
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

	return nil
}

func (s *Server) OnClose(c *gev.Connection) {
	log.Info("client connection closed %s", c.PeerAddr())
}

func (s *Server) Run() error {
	log.Info("server start")
	s.server.Start()
	return nil
}

func (s *Server) Init() error {
	h := tasks.NewTaskHandler(&s.config.Tasks)
	err := h.Init(s)
	if err != nil {
		return err
	}
	s.handlers = append(s.handlers, h)
	return nil
}

func (s *Server) ConfigDir() string {
	base_path := path.Dir(s.config_file)
	return base_path
}

func NewServer(cfg *config.ServerConf, cfgFile string) *Server {
	handler := &Server{
		config:      cfg,
		config_file: cfgFile,
	}
	bind_addr := cfg.Common.Bind + ":" + strconv.Itoa(cfg.Common.Port)
	log.Info("Bind on: %s", bind_addr)
	s, err := gev.NewServer(handler,
		gev.Network("tcp"),
		gev.Address(bind_addr),
		gev.NumLoops(cfg.Common.Loops),
		gev.CustomProtocol(protocol.New()))
	if err != nil {
		panic(err)
	}
	handler.server = s

	return handler
}
