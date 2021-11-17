package common

type Task interface {
	Start(c []byte) error
	Stop(c []byte) error
}

// Task Handler 接口
type TaskHandler interface {
	OnStart(*Task)
	OnError(*Task, error)
	OnStop(*Task, error)
	OnResult(config interface{}, result interface{}) error
}
