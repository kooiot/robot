package client

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kooiot/robot/client/config"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type Service struct {
	ctx    context.Context
	cancel context.CancelFunc

	config *config.ClientConf

	// Client
	client        *Client
	client_lock   sync.RWMutex
	connected_chn chan interface{}

	exit uint32 // 0 means not exit
}

func (s *Service) Run() error {
	client := NewClient(s.config, s.ctx)

	s.client_lock.Lock()
	s.client = client
	s.client_lock.Unlock()
	go client.Start(&s.connected_chn)

	go s.keepWorking()

	<-s.ctx.Done()
	return nil
}

func (s *Service) keepWorking() {
	xl := xlog.FromContextSafe(s.ctx)

	maxDelayTime := 20 * time.Second
	delayTime := time.Second

	for {
		select {
		case <-s.connected_chn:
			delayTime = time.Second
		case <-s.client.ClosedDoneChn():
			time.Sleep(delayTime)
			if atomic.LoadUint32(&s.exit) != 0 {
				return
			}

			delayTime = delayTime * 2
			if maxDelayTime < delayTime {
				delayTime = time.Second
			}

			xl.Info("try to reconnect to server...")
			client := NewClient(s.config, s.ctx)

			s.client_lock.Lock()
			s.client = client
			s.client_lock.Unlock()

			client.Start(&s.connected_chn)
		}
	}
}

func (s *Service) Close() {
	atomic.StoreUint32(&s.exit, 1)
	s.client_lock.Lock()
	s.client.Close()
	s.client_lock.Unlock()
	s.cancel()
}

func (s *Service) ConnectedChn() <-chan interface{} {
	return s.connected_chn
}

func NewService(cfg *config.ClientConf) *Service {
	ctx, cancel := context.WithCancel(context.Background())

	cli := &Service{
		config:        cfg,
		ctx:           xlog.NewContext(ctx, xlog.New()),
		cancel:        cancel,
		exit:          0,
		connected_chn: make(chan interface{}),
	}

	return cli
}
