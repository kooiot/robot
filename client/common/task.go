package common

import (
	"github.com/kooiot/robot/pkg/net/msg"
	uuid "github.com/satori/go.uuid"
)

type TaskResult struct {
	Result bool        `json:"result"`
	Error  string      `json:"error"`
	Info   interface{} `json:"info"`
}

type Task interface {
	ID() string
	Start() error
	Stop() error
}

type TaskBase struct {
	UUID string
	Info *msg.Task
}

func (t *TaskBase) ID() string {
	return t.UUID
}

func NewTaskBase(info *msg.Task) TaskBase {
	return TaskBase{
		UUID: uuid.NewV4().String(),
		Info: info,
	}
}

type TaskCreator func(handler TaskHandler, info *msg.Task, parent Task) Task
type TaskWait func(task Task, result TaskResult)

// Task Handler 接口
type TaskHandler interface {
	OnStart(Task)
	OnError(Task, error)
	OnSuccess(Task)
	OnResult(Task, TaskResult) error
	Spawn(creator TaskCreator, info *msg.Task, parent Task) Task
	Wait(Task, TaskWait) error
}
