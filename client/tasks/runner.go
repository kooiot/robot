package tasks

import (
	"errors"
	"sync"

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
	Info     *msg.Task  `json:"info"`
	Result   common.TaskResult
	Parent   common.Task
	Task     common.Task
	Children []TaskInfo
}

var gCreators = make(map[string]common.TaskCreator)

func RegisterTask(task_name string, creator common.TaskCreator) {
	gCreators[task_name] = creator
}

type Runner struct {
	tasks map[string]TaskInfo
	lock  sync.Mutex
}

func (r *Runner) OnStart(task common.Task) {
	r.lock.Lock()
	defer r.lock.Unlock()
	tinfo, err := r.tasks[task.ID()]
	if !err {
		tinfo.Status = ST_RUN
	}
}

func (r *Runner) OnError(task common.Task, err error) {
	result := common.TaskResult{
		Result: false,
		Error:  err.Error(),
	}
	r.OnResult(task, result)
}

func (r *Runner) OnSuccess(task common.Task) {
	result := common.TaskResult{
		Result: true,
		Error:  "done",
	}
	r.OnResult(task, result)
}

func (r *Runner) update_task_status(task common.Task) {
	// r.lock.Lock()
	// defer r.lock.Unlock()
	tinfo, ok := r.tasks[task.ID()]
	if !ok {
		log.Error("task not found!!")
	}

	all_done := true
	has_err := false
	for _, sub := range tinfo.Children {
		if sub.Status != ST_DONE && sub.Status != ST_FAILED {
			all_done = false
			break
		}
		if sub.Status == ST_FAILED {
			has_err = true
		}
	}
	if all_done {
		result := common.TaskResult{
			Result: true,
			Error:  "done",
		}
		if has_err {
			result.Result = false
			result.Error = "Sub task failed" // TODO:
		}
		//
		go r.OnResult(task, result)
	} else {
		log.Info("task not finished!")
	}
}

func (r *Runner) OnResult(task common.Task, result common.TaskResult) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	tinfo, ok := r.tasks[task.ID()]
	if !ok {
		log.Error("task not found!! result:%#v", result)
		return errors.New("task not found")
	}
	log.Info("task: %s result:%#v", tinfo.Info.Name, result)
	if result.Result {
		tinfo.Status = ST_DONE
	} else {
		tinfo.Status = ST_FAILED
	}
	tinfo.Result = result

	if tinfo.Parent != nil {
		r.update_task_status(tinfo.Parent)
	}

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

func (r *Runner) Spawn(creator common.TaskCreator, info *msg.Task, parent common.Task) common.Task {
	t := creator(r, info, parent)
	new_task := TaskInfo{
		Status: ST_NEW,
		Info:   info,
		Task:   t,
		Parent: parent,
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	r.tasks[t.ID()] = new_task
	if parent != nil {
		p_info := r.tasks[parent.ID()]
		p_info.Children = append(p_info.Children, new_task)
	}

	log.Info("Runner: spawn task:%s", new_task.Info.Name)
	go r.task_proc(t, &new_task)
	return t
}

func (r *Runner) Add(task *msg.Task, parent common.Task) (common.Task, error) {
	log.Info("Add task %s: %#v", task.Name, task.Option)
	creator := gCreators[task.Name]
	if creator == nil {
		log.Info("unknown task %s", task.Name)
		return nil, errors.New("unknown task " + task.Name)
	}
	t := r.Spawn(creator, task, parent)
	return t, nil
}

func (r *Runner) Wait(task common.Task, wait common.TaskWait) error {
	return nil
}

func NewRunner() *Runner {
	return &Runner{
		tasks: make(map[string]TaskInfo),
		lock:  sync.Mutex{},
	}
}
