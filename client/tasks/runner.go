package tasks

import (
	"context"
	"errors"
	"os/exec"
	"sync"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/client/config"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/log"
	"github.com/kooiot/robot/pkg/util/xlog"
	uuid "github.com/satori/go.uuid"
)

type TaskInfo struct {
	Info     *msg.Task       `json:"info"`
	Result   *msg.TaskResult `json:"result"`
	parent   common.Task
	task     common.Task
	children []*TaskInfo
	waits    []common.TaskWait
}

var gCreators = make(map[string]common.TaskCreator)

func RegisterTask(task_name string, creator common.TaskCreator) {
	if gCreators[task_name] != nil {
		log.Error("TaskName must be unique")
	}
	gCreators[task_name] = creator
}

type Runner struct {
	ctx      context.Context
	conf     *config.RunnerConf
	reporter common.Reporter
	tasks    map[string]TaskInfo
	lock     sync.Mutex
}

func (r *Runner) OnStart(task common.Task) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	tinfo, ok := r.tasks[task.TaskInfo().UUID]
	if ok {
		tinfo.Info.Status = msg.ST_RUN
		return r.reporter.SendTaskUpdate(task.TaskInfo())
	} else {
		return errors.New("task not found")
	}
}

func (r *Runner) OnError(task common.Task, err error) error {
	return r.OnResult(task, &msg.TaskResult{
		Result: false,
		Info:   err.Error(),
	})
}

func (r *Runner) OnSuccess(task common.Task) error {
	return r.OnResult(task, &msg.TaskResult{
		Result: true,
		Info:   "done",
	})
}

func (r *Runner) update_task_status(task common.Task) {
	xl := xlog.FromContextSafe(r.ctx)

	tinfo, ok := r.tasks[task.TaskInfo().UUID]
	if !ok {
		xl.Error("task not found!!")
		return
	}

	all_done := true
	has_err := false
	for _, sub := range tinfo.children {
		if sub.Info.Status != msg.ST_DONE && sub.Info.Status != msg.ST_FAILED {
			all_done = false
			break
		}
		if sub.Info.Status == msg.ST_FAILED {
			has_err = true
		}
	}

	if all_done {
		result := msg.TaskResult{
			Result: true,
			Info:   "done",
		}
		if has_err {
			result.Result = false
			result.Info = "Sub task failed" // TODO:
		}
		//
		xl.Info("parent task: %s status done !", tinfo.Info.Name)
		go r.OnResult(task, &result)
	} else {
		xl.Info("task not finished!")
	}
}

func (r *Runner) OnResult(task common.Task, result *msg.TaskResult) error {
	xl := xlog.FromContextSafe(r.ctx)

	r.lock.Lock()
	defer r.lock.Unlock()

	tinfo, ok := r.tasks[task.TaskInfo().UUID]
	if !ok {
		xl.Error("task not found!! result:%#v", result)
		return errors.New("task not found")
	}

	xl.Info("task: %s result:%#v", tinfo.Info.Name, result)

	result.Task = task.TaskInfo()
	if result.Result {
		tinfo.Info.Status = msg.ST_DONE
	} else {
		tinfo.Info.Status = msg.ST_FAILED
	}
	tinfo.Result = result

	if tinfo.parent != nil {
		xl.Info("update parent task status for %s!", tinfo.Info.Name)
		r.update_task_status(tinfo.parent)
	}

	for _, w := range tinfo.waits {
		w(task, result)
	}

	go r.ReportResult(&tinfo)

	return nil
}

func (r *Runner) task_proc(task common.Task, info *TaskInfo) {
	err := task.Start()
	if err != nil {
		r.OnError(task, err)
	} else {
		r.OnStart(task)
	}
}

func (r *Runner) Spawn(creator common.TaskCreator, info *msg.Task, parent common.Task) common.Task {
	xl := xlog.FromContextSafe(r.ctx)

	if len(info.UUID) == 0 {
		info.UUID = uuid.NewV4().String()
	}

	t := creator(r.ctx, r, info, parent)
	new_task := TaskInfo{
		Info:   info,
		task:   t,
		parent: parent,
	}
	info.Status = msg.ST_NEW

	r.lock.Lock()
	defer r.lock.Unlock()

	r.tasks[t.TaskInfo().UUID] = new_task
	if parent != nil {
		p_info := r.tasks[parent.TaskInfo().UUID]
		p_info.children = append(p_info.children, &new_task)
		info.ParentUUID = parent.TaskInfo().UUID
	}

	xl.Info("spawn task:%s", new_task.Info.Name)

	//
	r.reporter.SendTaskUpdate(t.TaskInfo())

	go r.task_proc(t, &new_task)

	return t
}

func (r *Runner) Add(task *msg.Task, parent common.Task) (common.Task, error) {
	xl := xlog.FromContextSafe(r.ctx)

	xl.Info("Add task %s: %#v", task.Name, task.Option)

	// Find Creator and create task
	creator := gCreators[task.Name]
	if creator == nil {
		xl.Error("unknown task %s", task.Name)
		return nil, errors.New("unknown task " + task.Name)
	}

	// Spawn Task
	t := r.Spawn(creator, task, parent)
	return t, nil
}

func (r *Runner) Wait(task common.Task, wait common.TaskWait) error {
	xl := xlog.FromContextSafe(r.ctx)
	xl.Info("Wait task %s", task.TaskInfo().Name)
	info := r.tasks[task.TaskInfo().UUID]
	if info.Info.Status != msg.ST_NEW && info.Info.Status != msg.ST_RUN {
		return errors.New("task is completed already")
	}

	info.waits = append(info.waits, wait)
	return nil
}

func (r *Runner) Halt() error {
	xl := xlog.FromContextSafe(r.ctx)
	if r.conf.Haltable {
		_, err := exec.Command("sh", "-c", "sleep 3 && halt").Output()
		if err != nil {
			xl.Error(err.Error())
			return err
		}
	}
	return nil
}

func (r *Runner) ReportResult(info *TaskInfo) error {
	r.reporter.SendResult(info.task.TaskInfo(), info.Result)
	return nil
}

func NewRunner(ctx context.Context, conf *config.RunnerConf, reporter common.Reporter) *Runner {
	xl := xlog.FromContextSafe(ctx).Spawn().AppendPrefix("Runner")
	return &Runner{
		ctx:      xlog.NewContext(ctx, xl),
		conf:     conf,
		reporter: reporter,
		tasks:    make(map[string]TaskInfo),
		lock:     sync.Mutex{},
	}
}
