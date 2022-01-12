package common

import (
	"context"

	"github.com/kooiot/robot/pkg/net/msg"
)

type Task interface {
	TaskInfo() *msg.Task
	Start() error
	Stop() error
}

type TaskBase struct {
	Info *msg.Task `json:"info"`
}

func (t *TaskBase) TaskInfo() *msg.Task {
	return t.Info
}

func NewTaskBase(info *msg.Task) TaskBase {
	return TaskBase{
		Info: info,
	}
}

type TaskCreator func(ctx context.Context, handler TaskHandler, info *msg.Task, parent Task) Task
type TaskWait func(task Task, result *msg.TaskResult)

// Task Handler 接口
type TaskHandler interface {
	OnStart(Task)
	OnError(Task, error)
	OnSuccess(Task)
	OnResult(Task, *msg.TaskResult) error
	Spawn(creator TaskCreator, info *msg.Task, parent Task) Task
	Wait(Task, TaskWait) error
}

type Reporter interface {
	SendResult(*msg.Task, *msg.TaskResult) error
	SendTaskUpdate(*msg.Task) error
}
