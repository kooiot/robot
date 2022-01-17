package tasks

import (
	"context"
	"errors"
	"os/exec"
	"sync"
	"time"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/client/config"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/log"
	"github.com/kooiot/robot/pkg/util/xlog"
	uuid "github.com/satori/go.uuid"
)

type TaskInfo struct {
	Info     msg.Task       `json:"info"`
	Result   msg.TaskResult `json:"result"`
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
	tasks    map[string]*TaskInfo
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
	return r.OnResult(task, msg.TaskResult{
		Result: false,
		Info:   err.Error(),
	})
}

func (r *Runner) OnSuccess(task common.Task) error {
	return r.OnResult(task, msg.TaskResult{
		Result: true,
		Info:   "done",
	})
}

func (r *Runner) get_task_status(task_id string) (msg.StatusType, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	tinfo, ok := r.tasks[task_id]
	if ok {
		return tinfo.Info.Status, nil
	} else {
		return -1, errors.New("not found!")
	}
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
		xl.Debug("parent task: %s status done !", tinfo.Info.Task)
		go r.OnResult(task, result)
	} else {
		xl.Debug("task not finished!")
	}
}

func (r *Runner) OnResult(task common.Task, result msg.TaskResult) error {
	xl := xlog.FromContextSafe(r.ctx)

	r.lock.Lock()
	defer r.lock.Unlock()

	tinfo, ok := r.tasks[task.TaskInfo().UUID]
	if !ok {
		xl.Error("task not found!! result:%#v", result)
		return errors.New("task not found")
	}

	result.Task = task.TaskInfo()
	xl.Info("task: %s result:%#v", tinfo.Info.Task, result)

	if result.Result {
		tinfo.Info.Status = msg.ST_DONE
	} else {
		tinfo.Info.Status = msg.ST_FAILED
	}
	tinfo.Result = result

	if tinfo.parent != nil {
		xl.Debug("update parent task status for %s!", tinfo.Info.Task)
		r.update_task_status(tinfo.parent)
	}

	for _, w := range tinfo.waits {
		w(task, result)
	}

	go r.ReportResult(tinfo)

	return nil
}

func (r *Runner) task_proc(task common.Task, info *TaskInfo) {
	xl := xlog.FromContextSafe(r.ctx)
	if len(info.Info.Depends) > 0 {
		for {
			xl.Error("task: %s depends size: %d", info.Info.ID, len(info.Info.Depends))
			xl.Info("task %s depends ......... %#v", info.Info.ID, task)
			time.Sleep(1 * time.Second)
			xl.Info("task %s wait depends task finish......... %#v", info.Info.ID, info.Info.Depends)
			xl.Info("task %s wait depends task finish......... %#v", info.Info.ID, task)
			all_done := true
			for _, dep := range info.Info.Depends {
				sts, err := r.get_task_status(dep)
				if err != nil {
					// Task does not exists
					r.OnError(task, errors.New("depends task "+dep+" not exists"))
					return
				}
				if sts == msg.ST_FAILED {
					r.OnError(task, errors.New("depends task "+dep+" failed"))
					return
				}
				if sts != msg.ST_DONE {
					all_done = false
					break
				}
			}
			if all_done {
				xl.Info("all depends task finished!")
				break
			}
		}
	} else {
		xl.Error("task: %s depends null", info.Info.ID)
	}

	err := task.Start()
	if err != nil {
		r.OnError(task, err)
	} else {
		r.OnStart(task)
	}
}

func (r *Runner) Spawn(creator common.TaskCreator, info msg.Task, parent common.Task) common.Task {
	xl := xlog.FromContextSafe(r.ctx)

	if len(info.UUID) == 0 {
		info.UUID = uuid.NewV4().String()
	}

	t := creator(r.ctx, r, info, parent)
	new_task := &TaskInfo{
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
		p_info.children = append(p_info.children, new_task)
		info.ParentUUID = parent.TaskInfo().UUID
	}

	xl.Debug("spawn task: %s", new_task.Info.ID)

	//
	r.reporter.SendTaskUpdate(t.TaskInfo())

	go r.task_proc(t, new_task)

	return t
}

func (r *Runner) Add(task msg.Task, parent common.Task) (common.Task, error) {
	xl := xlog.FromContextSafe(r.ctx)

	xl.Debug("Add task %s: %#v", task.Task, task.Option)

	// Find Creator and create task
	creator := gCreators[task.Task]
	if creator == nil {
		xl.Error("unknown task %s", task.Task)
		return nil, errors.New("unknown task " + task.Task)
	}

	// Spawn Task
	t := r.Spawn(creator, task, parent)
	return t, nil
}

func (r *Runner) Wait(task common.Task, wait common.TaskWait) error {
	xl := xlog.FromContextSafe(r.ctx)
	xl.Debug("Wait task %s", task.TaskInfo().Task)
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
			xl.Error("System halt error: %s", err.Error())
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
		tasks:    make(map[string]*TaskInfo),
		lock:     sync.Mutex{},
	}
}
