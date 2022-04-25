package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kooiot/robot/server/common"
)

type RobotRouter struct {
}

var robotRouter RobotRouter

func (b *RobotRouter) GetStats(c *gin.Context) {
	stats := common.RobotStats{
		UpdateTime:   time.Now().String(),
		RunToday:     []common.RunResult{},
		ClientActive: []common.RunResult{},
		OrderSource:  []common.RunResult{},
		ClientLevel:  []common.RunResult{},
		ErrorTop:     []common.ErrorData{},
		TestStatus:   common.TestStatus{},
		ServerStatus: common.ServerStatus{},
	}
	OkWithData(stats, c)
}

func (s *RobotRouter) InitRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	router := Router.Group("robot")
	{
		router.POST("stats", robotRouter.GetStats)
	}
	return router
}
