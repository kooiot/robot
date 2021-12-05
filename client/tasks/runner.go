package tasks

import (
	"errors"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/log"
)

type StatusType int32

const (
	ST_NEW    StatusType = 0
	ST_RUN    StatusType = 1
	ST_FAILED StatusType = 2
	ST_DONE   StatusType = 4
	// ST_SPAWN  StatusType = 8
)

type TaskInfo struct {
	Status   StatusType `json:"status"`
	Info     *msg.Task  `json"info"`
	Parent   common.Task
	Task     common.Task
	Children []TaskInfo
}

var gCreators = make(map[string]common.TaskCreator)

func RegisterTask(task_name string, creator common.TaskCreator) {
	gCreators[task_name] = creator
}

type Runner struct {
	tasks map[common.Task]TaskInfo
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

func (r *Runner) task_proc(task common.Task, info *TaskInfo) {
	log.Info("Runner: start task:%s", info.Info.Name)
	err := task.Start()
	if err != nil {
		log.Error("Runner: start task:%s error: %s", info.Info.Name, err.Error())
		r.OnError(task, err)
	} else {
		r.OnStart(task)
	}
}

func (r *Runner) Spawn(creator common.TaskCreator, info *msg.Task, parent common.Task) {
	t := creator(r, info)
	new_task := TaskInfo{
		Status: ST_NEW,
		Info:   info,
		Task:   t,
		Parent: parent,
	}
	r.tasks[t] = new_task
	p_info := r.tasks[parent]
	p_info.Children = append(p_info.Children, new_task)

	go r.task_proc(t, &new_task)
}

func (r *Runner) Add(task *msg.Task, parent common.Task) error {
	log.Info("%s: %v", task.Name, task.Option)
	creator := gCreators[task.Name]
	if creator == nil {
		log.Info("unknown task %s", task.Name)
		return errors.New("unknown task " + task.Name)
	}
	r.Spawn(creator, task, nil)
	return nil
}

func NewRunner() *Runner {
	return &Runner{
		tasks: make(map[common.Task]TaskInfo),
	}
}
