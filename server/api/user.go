package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kooiot/robot/server/common"
	uuid "github.com/satori/go.uuid"
)

type UserRouter struct {
}

var userRouter UserRouter

// @Tags Base
// @Summary 用户登录
// @Produce  application/json
// @Param data body systemReq.Login true "用户名, 密码, 验证码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /base/login [post]
func (b *UserRouter) Login(c *gin.Context) {
	req := common.LoginReq{}
	_ = c.ShouldBindJSON(&req)

	if req.Password == "admin" && req.Username == "admin" {
		resp := common.LoginResp{
			Username: "admin",
			Name:     "Admin",
			UUID:     uuid.NewV4().String(),
			Token:    "adminToken",
		}
		OkWithData(resp, c)
	} else {
		FailWithMessage("验证码错误", c)
	}
}

func (b *UserRouter) GetUserInfo(c *gin.Context) {
	token := c.GetHeader("authorization")

	if token == "adminToken" {
		resp := common.UserInfo{
			UserID:   1,
			Username: "admin",
			RealName: "Admin",
			Desc:     "Admin",
			Avatar:   "https://q1.qlogo.cn/g?b=qq&nk=190848757&s=640",
			HomePath: "/dashboard/analysis",
		}
		resp.Roles = append(resp.Roles, common.UserRole{
			RoleName: "Super Admin",
			Value:    "super",
		})
		OkWithData(resp, c)
	} else {
		FailWithMessage("验证码错误", c)
	}
}

func (b *UserRouter) GetPermCode(c *gin.Context) {
	token := c.GetHeader("authorization")

	if token == "adminToken" {
		resp := []string{}
		resp = append(resp, "1000")
		resp = append(resp, "3000")
		resp = append(resp, "5000")
		OkWithData(resp, c)
	} else {
		FailWithMessage("验证码错误", c)
	}
}

func (b *UserRouter) Logout(c *gin.Context) {
	token := c.GetHeader("authorization")

	if token == "adminToken" {
		OkWithMessage("Token has been destroyed", c)
	} else {
		FailWithMessage("验证码错误", c)
	}
}

func (s *UserRouter) InitUserRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	router := Router.Group("api")
	{
		router.GET("getUserInfo", userRouter.GetUserInfo)
		router.GET("getPermCode", userRouter.GetPermCode)
	}
	return router
}
