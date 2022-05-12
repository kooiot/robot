package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kooiot/robot/server/api/middleware"
	"github.com/kooiot/robot/server/config"
)

func IndexToAdmin(c *gin.Context) {
	c.Redirect(http.StatusFound, "/admin")
}

// 初始化总路由

func Routers(conf *config.HttpApiConf) *gin.Engine {
	var Router = gin.Default()
	Router.StaticFS("/admin/", http.Dir(conf.Static)) // 静态地址

	// "github.com/gin-contrib/static"
	// Router.Use(static.Serve("/", static.LocalFile(conf.Static, false)))
	// 方便统一添加路由组前缀 多服务器上线使用

	Router.GET("/", IndexToAdmin)
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
