package server

import (
	"context"
	"encoding/json"
	"path"
	"strconv"
	"time"

	"github.com/Allenxuxu/gev"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/net/protocol"
	"github.com/kooiot/robot/pkg/util/xlog"
	"github.com/kooiot/robot/server/common"
	"github.com/kooiot/robot/server/config"
	"github.com/kooiot/robot/server/tasks"
)

type Server struct {
	ctx         context.Context
	cancel      context.CancelFunc
	config      *config.ServerConf
	config_file string
	server      *gev.Server
	handlers    []common.ServerHandler
	clients     []common.Client
}

func (s *Server) OnConnect(c *gev.Connection) {
	xl := xlog.FromContextSafe(s.ctx)
	xl.Info("client connected: %s", c.PeerAddr())
}

func (s *Server) AfterLogin(conn *gev.Connection, client *common.Client) {
	xl := xlog.FromContextSafe(s.ctx)
	time.Sleep(1 * time.Second)

	xl.Debug("AfterLogin %s", conn.PeerAddr())
	for _, h := range s.handlers {
		h.AfterLogin(conn, client)
	}
}

func (s *Server) HandleLogin(c *gev.Connection, req *msg.Login) interface{} {
	xl := xlog.FromContextSafe(s.ctx)
	xl.Trace("received login: %v", req)

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
		xl.Error("failed encode resp: %s", err)
		return nil
	} else {
		go s.AfterLogin(c, &client)
		return protocol.PackMessage("login_resp", data)
	}
}

func (s *Server) HandleHeartbeat(c *gev.Connection, req *msg.HeartBeat) interface{} {
	xl := xlog.FromContextSafe(s.ctx)
	xl.Info("received heartbeat: %#v", req)

	resp := &msg.HeartBeat{
		ID:   req.ID,
		Time: time.Now().UTC().Unix(),
	}
	// xl.Debug("Send back: %v", resp)

	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return protocol.PackMessage("heartbeat", data)
}

func (s *Server) HandleTaskUpdate(c *gev.Connection, req *msg.Task) interface{} {
	xl := xlog.FromContextSafe(s.ctx)
	xl.Info("received task: %#v", req)
	return nil
}

func (s *Server) HandleTaskResult(c *gev.Connection, req *msg.TaskResult) interface{} {
	xl := xlog.FromContextSafe(s.ctx)
	xl.Info("received result: %#v", req)
	return nil
}

func (s *Server) OnMessage(c *gev.Connection, ctx interface{}, data []byte) interface{} {
	xl := xlog.FromContextSafe(s.ctx)
	msgType := ctx.(string)

	switch msgType {
	case "login":
		req := msg.Login{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Error(err.Error())
		}
		return s.HandleLogin(c, &req)
	case "logout":
		req := msg.Logout{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Error(err.Error())
		}

		resp := &msg.Response{
			Content: "OK",
		}

		data, err := json.Marshal(resp)
		if err != nil {
			xl.Error("failed encode resp: %s", err)
		} else {
			return protocol.PackMessage("logout_resp", data)
		}
	case "heartbeat":
		req := msg.HeartBeat{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Info(err.Error())
		}
		return s.HandleHeartbeat(c, &req)
	case "task.update":
		req := msg.Task{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Error(err.Error())
		}
		return s.HandleTaskUpdate(c, &req)
	case "task.result":
		req := msg.TaskResult{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Error(err.Error())
		}
		return s.HandleTaskResult(c, &req)
	default:
		xl.Error("unknown msg type %s", msgType)
	}

	return nil
}

func (s *Server) OnClose(c *gev.Connection) {
	xl := xlog.FromContextSafe(s.ctx)
	xl.Info("client connection closed %s", c.PeerAddr())
}

func (s *Server) Run() error {
	xl := xlog.FromContextSafe(s.ctx)
	xl.Info("server start")
	s.server.Start()
	return nil
}

func (s *Server) Init() error {
	h := tasks.NewTaskHandler(s.ctx, &s.config.Tasks)
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
	ctx, cancel := context.WithCancel(context.Background())
	handler := &Server{
		ctx:         xlog.NewContext(ctx, xlog.New()),
		cancel:      cancel,
		config:      cfg,
		config_file: cfgFile,
	}
	xl := xlog.FromContextSafe(handler.ctx)
	bind_addr := cfg.Common.Bind + ":" + strconv.Itoa(cfg.Common.Port)
	xl.Info("Bind on: %s", bind_addr)
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
