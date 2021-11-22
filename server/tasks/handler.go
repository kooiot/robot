package tasks

import (
	"github.com/kooiot/robot/server/common"
	"github.com/kooiot/robot/server/config"
)

type TaskHandler struct {
	autos *config.AutoTasks
}

func (h *TaskHandler) Init(server common.Server) error {
	return nil
}

func (h *TaskHandler) AfterLogin(client *common.Client) {

}

func (h *TaskHandler) BeforeLogout(client *common.Client) {

}

func NewTaskHandler(autos *config.AutoTasks) *TaskHandler {
	return &TaskHandler{autos: autos}
}
