package api

import (
	"github.com/gin-gonic/gin"
)

type BaseRouter struct {
}

var baseRouter BaseRouter

// @Tags Base
// @Summary 用户登录
// @Produce  application/json
// @Param data body systemReq.Login true "用户名, 密码, 验证码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /base/login [post]
func (b *BaseRouter) Login(c *gin.Context) {
	//_ = c.ShouldBindJSON(&l)

	FailWithMessage("验证码错误", c)
}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	router := Router.Group("base")
	{
		router.POST("login", baseRouter.Login)
	}
	return router
}
