package tasks

import (
	"encoding/json"
	"errors"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/pkg/net/msg"
)

type BatchTask struct {
	config  msg.BatchTask
	handler common.TaskHandler
}

func init() {
	RegisterTask("batch", NewBatchTask)
}

func (s *BatchTask) Start() error {
	r, ok := s.handler.(*Runner)
	if !ok {
		return errors.New("error object")
	}
	for _, t := range s.config.Tasks {
		r.Add(&t, s)
	}
	return nil
}

func (s *BatchTask) Stop() error {
	return nil
}

func NewBatchTask(handler common.TaskHandler, option interface{}) common.Task {
	data, _ := json.Marshal(option)

	conf := msg.BatchTask{}
	json.Unmarshal(data, &conf)

	return &BatchTask{
		config: conf,
	}
}
