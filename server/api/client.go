package api

import (
	"github.com/gin-gonic/gin"
)

type ClientRouter struct {
}

var clientRouter ClientRouter

// @Tags Client
// @Summary 枚举客户端列表
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /base/login [post]
func (b *ClientRouter) List(c *gin.Context) {
	//_ = c.ShouldBindJSON(&l)

	FailWithMessage("错误", c)
}

func (s *ClientRouter) InitClientRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	router := Router.Group("client")
	{
		router.POST("list", clientRouter.List)
	}
	return router
}
