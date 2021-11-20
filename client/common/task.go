package common

import "github.com/kooiot/robot/pkg/net/msg"

type Task interface {
	Start() error
	Stop() error
}

type TaskCreator func(TaskHandler, interface{}) Task

// Task Handler 接口
type TaskHandler interface {
	OnStart(Task)
	OnError(Task, error)
	OnStop(Task, error)
	OnResult(config interface{}, result interface{}) error
	Spawn(creator TaskCreator, info *msg.Task, parent Task)
}

type TaskResult struct {
	Result bool   `json:"result"`
	Error  string `json:"error"`
}
