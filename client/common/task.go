package common

import (
	"context"
	"errors"

	"github.com/kooiot/robot/pkg/net/msg"
)

type Task interface {
	ID() string
	UUID() string
	Init() error
	Run() (interface{}, error)
	Abort() error
}

type TaskBase struct {
	Info msg.Task `json:"info"`
}

func (t *TaskBase) ID() string {
	return t.Info.ID
}

func (t *TaskBase) UUID() string {
	return t.Info.UUID
}

func (t *TaskBase) Init() error {
	return nil
}

func (t *TaskBase) Abort() error {
	return errors.New("abort not implemented")
}

func NewTaskBase(info msg.Task) TaskBase {
	return TaskBase{
		Info: info,
	}
}

type TaskCreator func(ctx context.Context, handler TaskHandler, info msg.Task, parent Task) Task
type TaskWait func(task Task, result msg.TaskResultDetail)

// Task Handler 接口
type TaskHandler interface {
	Spawn(creator TaskCreator, info msg.Task, parent Task) Task
	Wait(Task, TaskWait) error
}

type Reporter interface {
	PrintDone()
	SendResult(*msg.TaskResult) error
	SendTaskUpdate(*msg.Task) error
}
