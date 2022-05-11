package api

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/kooiot/robot/server/api/middleware"
	"github.com/kooiot/robot/server/config"
)

// 初始化总路由

func Routers(conf *config.HttpApiConf) *gin.Engine {
	var Router = gin.Default()
	// Router.StaticFS("/", http.Dir(conf.Static)) // 静态地址
	Router.Use(static.Serve("/", static.LocalFile("static", false)))
	// 方便统一添加路由组前缀 多服务器上线使用

	PublicGroup := Router.Group("api")
	{
		baseRouter.InitPublicRouter(PublicGroup)
	}
	PrivateGroup := Router.Group("api")
	PrivateGroup.Use(middleware.BackendTokenAuth())
	{
		baseRouter.InitRouter(PublicGroup)
		menuRouter.InitRouter(PrivateGroup)
		robotRouter.InitRouter(PrivateGroup)
	}

	return Router
}
