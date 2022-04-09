package api

import (
	"fmt"
	"time"

	"github.com/fvbock/endless"
	"github.com/kooiot/robot/server/config"
)

func RunServer(cfg *config.HttpApiConf) error {
	router := Routers(cfg)
	address := fmt.Sprintf("%s:%d", cfg.Bind, cfg.Port)

	srv := endless.NewServer(address, router)
	srv.ReadHeaderTimeout = 10 * time.Millisecond
	srv.WriteTimeout = 10 * time.Second
	srv.MaxHeaderBytes = 1 << 20

	return srv.ListenAndServe()
}
