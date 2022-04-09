package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kooiot/robot/server/config"
)

// 初始化总路由

func Routers(conf *config.HttpApiConf) *gin.Engine {
	var Router = gin.Default()
	Router.StaticFS(conf.Static, http.Dir(conf.Static)) // 静态地址
	// 方便统一添加路由组前缀 多服务器上线使用

	PublicGroup := Router.Group("")
	{
		baseRouter.InitBaseRouter(PublicGroup)
	}

	return Router
}
