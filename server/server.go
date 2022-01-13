package server

import (
	"context"
	"encoding/json"
	"path"
	"strconv"
	"sync"
	"sync/atomic"
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

	client_next_id int32
	clients        map[int32]*common.Client
	clients_lock   sync.RWMutex

	// write to this channel to write the message sent to server
	send_chn chan (*msg.Message)
	// read from this channel to get the next message sent by server
	read_chn chan (*msg.Message)

	exit uint32 // 0 means not exit
}

func (s *Server) OnConnect(c *gev.Connection) {
	xl := xlog.FromContextSafe(s.ctx)
	xl.Info("client connected: %s", c.PeerAddr())
}

func (s *Server) after_login(conn *gev.Connection, client *common.Client) {
	xl := xlog.FromContextSafe(s.ctx)
	time.Sleep(1 * time.Second)

	xl.Debug("AfterLogin %s", conn.PeerAddr())
	for _, h := range s.handlers {
		h.AfterLogin(conn, client)
	}
}

func (s *Server) handle_login(c *gev.Connection, req *msg.Login) interface{} {
	xl := xlog.FromContextSafe(s.ctx)
	xl.Trace("received login: %#v", req)

	var client_id int32
	for {
		client_id = atomic.AddInt32(&s.client_next_id, 1)
		if cli := s.get_client_by_id(client_id); cli == nil {
			break
		}
	}

	resp := msg.LoginResp{
		ClientID: req.ClientID,
		ID:       client_id,
		Reason:   "OK",
	}
	client := common.NewClient(client_id, c, req)

	s.clients_lock.Lock()
	s.clients[client_id] = client
	s.clients_lock.Unlock()

	data, err := json.Marshal(resp)
	if err != nil {
		xl.Error("failed encode resp: %s", err)
		return nil
	} else {
		go s.after_login(c, client)
		return protocol.PackMessage("login_resp", data)
	}
}

func (s *Server) get_client_by_id(id int32) *common.Client {
	s.clients_lock.Lock()
	defer s.clients_lock.Unlock()
	return s.clients[id]
}

func (s *Server) get_client_by_conn(conn *gev.Connection) *common.Client {
	s.clients_lock.Lock()
	defer s.clients_lock.Unlock()
	for _, c := range s.clients {
		if c.Conn == conn {
			return c
		}
	}
	return nil
}

func (s *Server) remove_client(cli *common.Client) {
	s.clients_lock.Lock()
	defer s.clients_lock.Unlock()
	delete(s.clients, cli.ID)
}

func (s *Server) check_heartbeat() {
	s.clients_lock.Lock()
	now := time.Now()
	var to_closed []*common.Client
	for _, c := range s.clients {
		if now.Sub(c.LastHeartbeat) > 90*time.Second {
			to_closed = append(to_closed, c)
		}
	}
	s.clients_lock.Unlock()

	for _, c := range to_closed {
		c.Conn.Close()
	}
}

func (s *Server) handle_logout(c *gev.Connection, req *msg.Logout) interface{} {
	xl := xlog.FromContextSafe(s.ctx)
	cli := s.get_client_by_id(req.ID)
	if cli == nil {
		return nil
	} else {
		s.remove_client(cli)
		go func() {
			time.Sleep(1 * time.Second)
			cli.Conn.Close()
		}()
	}

	resp := &msg.Response{
		Content: "OK",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		xl.Error("failed encode resp: %s", err)
		return nil
	} else {
		return protocol.PackMessage("logout_resp", data)
	}
}

func (s *Server) handle_heartbeat(c *gev.Connection, req *msg.HeartBeat) interface{} {
	xl := xlog.FromContextSafe(s.ctx)
	xl.Info("received heartbeat: %#v", req)

	resp := &msg.HeartBeat{
		ID:   req.ID,
		Time: time.Now().UTC().Unix(),
	}
	// xl.Debug("Send back: %v", resp)

	cli := s.get_client_by_id(req.ID)
	if cli != nil {
		cli.LastHeartbeat = time.Now()
	} else {
		// Remove client
		c.Close()
		return nil
	}

	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return protocol.PackMessage("heartbeat", data)
}

func (s *Server) handle_task_update(c *gev.Connection, req *msg.Task) interface{} {
	xl := xlog.FromContextSafe(s.ctx)
	xl.Info("received task: %#v", req)
	return nil
}

func (s *Server) handle_task_result(c *gev.Connection, req *msg.TaskResult) interface{} {
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
			xl.Error("JSON.Unmarshal error: %s", err.Error())
		}
		return s.handle_login(c, &req)
	case "logout":
		req := msg.Logout{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Error("JSON.Unmarshal error: %s", err.Error())
		}
		return s.handle_logout(c, &req)
	case "heartbeat":
		req := msg.HeartBeat{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Info("JSON.Unmarshal error: %s", err.Error())
		}
		return s.handle_heartbeat(c, &req)
	case "task.update":
		req := msg.Task{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Error("JSON.Unmarshal error: %s", err.Error())
		}
		return s.handle_task_update(c, &req)
	case "task.result":
		req := msg.TaskResult{}
		if err := json.Unmarshal(data, &req); err != nil {
			xl.Error("JSON.Unmarshal error: %s", err.Error())
		}
		return s.handle_task_result(c, &req)
	default:
		xl.Error("unknown msg type %s", msgType)
	}

	return nil
}

func (s *Server) OnClose(c *gev.Connection) {
	xl := xlog.FromContextSafe(s.ctx)
	xl.Info("client connection closed %s", c.PeerAddr())
	cli := s.get_client_by_conn(c)
	if cli != nil {
		s.remove_client(cli)
	}
}

func (s *Server) Run() error {
	xl := xlog.FromContextSafe(s.ctx)
	go s.worker()

	xl.Info("server start")
	s.server.Start() // blocked here
	return nil
}

func (s *Server) worker() {
	hbCheck := time.NewTicker(time.Second)
	defer hbCheck.Stop()

	for {
		select {
		case <-hbCheck.C:
			if atomic.LoadUint32(&s.exit) != 0 {
				return
			}
			s.check_heartbeat()
		}
	}
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

func (s *Server) Close() {
	atomic.StoreUint32(&s.exit, 1)
	s.cancel()
}

func (s *Server) ConfigDir() string {
	base_path := path.Dir(s.config_file)
	return base_path
}

func NewServer(cfg *config.ServerConf, cfgFile string) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	handler := &Server{
		ctx:            xlog.NewContext(ctx, xlog.New()),
		cancel:         cancel,
		config:         cfg,
		config_file:    cfgFile,
		send_chn:       make(chan *msg.Message, 100),
		read_chn:       make(chan *msg.Message, 100),
		client_next_id: 0,
		clients:        make(map[int32]*common.Client),
		clients_lock:   sync.RWMutex{},
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
