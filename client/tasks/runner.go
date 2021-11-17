package tasks

import (
	"github.com/kooiot/robot/client/common"
)

type Runner struct {
	tasks []common.Task
}

func (r *Runner) OnStart(common.Task) {

}

func (r *Runner) OnError(common.Task, error) {

}

func (r *Runner) OnStop(common.Task, error) {

}

func (r *Runner) OnResult(config interface{}, result interface{}) error {

	return nil
}

func (r *Runner) Spawn(creator func(common.TaskHandler, interface{}) common.Task, option interface{}) {
	t := creator(r, option)
	r.tasks = append(r.tasks, t)
}

func NewRunner() *Runner {
	return &Runner{}
}
