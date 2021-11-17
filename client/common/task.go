package common

type Task interface {
	Start() error
	Stop() error
}

// Task Handler 接口
type TaskHandler interface {
	OnStart(Task)
	OnError(Task, error)
	OnStop(Task, error)
	OnResult(config interface{}, result interface{}) error
	Spawn(creator func(TaskHandler, interface{}) Task, option interface{})
}

type TaskResult struct {
	Result bool   `json:"result"`
	Error  string `json:"error"`
}
