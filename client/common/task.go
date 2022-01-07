package common

import "github.com/kooiot/robot/pkg/net/msg"

type TaskResult struct {
	Result bool        `json:"result"`
	Error  string      `json:"error"`
	Info   interface{} `json:"info"`
}

type Task interface {
	Start() error
	Stop() error
}

type TaskCreator func(TaskHandler, *msg.Task) Task

// Task Handler 接口
type TaskHandler interface {
	OnStart(Task)
	OnError(Task, error)
	OnSuccess(Task)
	OnResult(Task, TaskResult) error
	Spawn(creator TaskCreator, info *msg.Task, parent Task) Task
}
