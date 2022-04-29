package api

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/server/common"
)

type RobotRouter struct {
}

var robotRouter RobotRouter

func (b *RobotRouter) GenStats(stats *common.RobotStats) {
	for i := 0; i < 24; i++ {
		success := rand.Intn(60)
		fail := rand.Intn(success + 1)
		stats.RunToday = append(stats.RunToday, common.RunResult{Hour: strconv.Itoa(i), Success: uint32(success), Fail: uint32(fail)})
	}

	for i := 0; i < 10; i++ {
		stats.ErrorTop = append(stats.ErrorTop, common.ErrorData{ID: uint32(i), Name: strconv.Itoa(i) + "Test result", Count: uint32(20 - i)})
	}

	for i := 0; i < 5; i++ {
		stats.Clients = append(stats.Clients, common.ClientData{ID: int32(i + 100), Info: msg.Login{
			ClientID: "AAAAA" + strconv.Itoa(i),
		}, Online: time.Now().Format("2006-01-02 15:04:05"), Status: "Running"})
	}

	for i := 0; i < 5; i++ {
		stats.ErrorStats = append(stats.ErrorStats, common.ErrorStatsData{Name: "AAAAA" + strconv.Itoa(i), Count: uint32(i)})
	}
}

func (b *RobotRouter) GetStats(c *gin.Context) {
	// stats := common.RobotStats{
	// 	UpdateTime:   time.Now().Format("2006-01-02 15:04:05"),
	// 	RunToday:     []common.RunResult{},
	// 	ClientActive: []common.RunResult{},
	// 	ErrorStats:   []common.ErrorStatsData{},
	// 	ClientLevel:  []common.RunResult{},
	// 	ErrorTop:     []common.ErrorData{},
	// 	TestStatus:   common.TestStatus{},
	// 	ServerStatus: common.ServerStatus{},
	// }
	// b.GenStats(&stats)
	stats := G_Stats.GetStats()
	OkWithData(stats, c)
}

func (b *RobotRouter) GetDetail(c *gin.Context) {
	req := common.DetailQuery{}
	_ = c.ShouldBindJSON(&req)
	if req.ID == 0 {
		FailWithMessage("unknown client id", c)
		return
	}

	detail := G_Stats.GetDetail(req.ID)

	OkWithData(detail, c)
}

func (b *RobotRouter) GetInfo(c *gin.Context) {
	req := common.InfoQuery{}
	_ = c.ShouldBindJSON(&req)
	if req.ClientID == "" {
		FailWithMessage("unknown client id", c)
		return
	}

	data := G_Stats.GetInfo(req.ClientID)

	OkWithData(data, c)
}

func (s *RobotRouter) InitRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	router := Router.Group("robot")
	{
		router.POST("stats", robotRouter.GetStats)
		router.POST("detail", robotRouter.GetDetail)
		router.POST("info", robotRouter.GetInfo)
	}
	return router
}
