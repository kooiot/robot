package pong

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kooiot/robot/pkg/util/log"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type PongServerConf struct {
	Log   log.LogConf    `mapstructure:"log" json:"log"`
	Ports []SerialConfig `mapstructure:"ports" json:"ports"`
}

// GetDefaultClientConf returns a client configuration with default values.
func GetDefaultServerConf() PongServerConf {
	return PongServerConf{
		Log: log.LogConf{
			Filename: "pong.log",
			Dir:      "log",
			Level:    "info",
		},
	}
}

func ParseServerConfig(path string) (PongServerConf, error) {
	cfg := GetDefaultServerConf()
	err := cfg.Load(path)
	return cfg, err
}

func (cfg *PongServerConf) Complete() {
	// fmt.Printf("ProxyL: %v\n", cfg.Proxy)

	// if cfg.LogLink == "console" {
	// 	cfg.LogDir = "console"
	// } else {
	// 	cfg.LogDir = "file"
	// }
}

func (cfg *PongServerConf) Validate() error {
	return nil
}

func (cfg *PongServerConf) Load(path ...string) error {
	v := Viper(cfg, path...)
	if v == nil {
		return fmt.Errorf("invalid protocol")
	}
	cfg.Complete()

	return nil
}

type Service struct {
	ctx    context.Context
	cancel context.CancelFunc

	config *PongServerConf

	// Client
	servers []*PongServer
	lock    sync.RWMutex

	exit uint32 // 0 means not exit
}

func (s *Service) Run() error {

	s.lock.Lock()
	for _, port := range s.config.Ports {
		server := NewPongServer(s.ctx, port)
		s.servers = append(s.servers, server)
		go server.Run()
	}
	s.lock.Unlock()

	go s.keepWorking()

	<-s.ctx.Done()
	return nil
}

func (s *Service) keepWorking() {
	for {
		if atomic.LoadUint32(&s.exit) != 0 {
			return
		}
		s.lock.RLock()
		for _, s := range s.servers {
			if !s.IsRunning() {
				panic("port failure")
			}
		}
		s.lock.RUnlock()
		time.Sleep(time.Second)
	}
}

func (s *Service) Close() {
	atomic.StoreUint32(&s.exit, 1)
	s.lock.Lock()
	for _, s := range s.servers {
		s.Abort()
	}
	s.lock.Unlock()
	s.cancel()
}

func NewService(cfg *PongServerConf) *Service {
	ctx, cancel := context.WithCancel(context.Background())

	cli := &Service{
		config: cfg,
		ctx:    xlog.NewContext(ctx, xlog.New()),
		cancel: cancel,
		exit:   0,
	}

	return cli
}
