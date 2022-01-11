package tasks

import (
	"encoding/json"
	"errors"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/log"
)

type BatchTask struct {
	common.TaskBase
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
		log.Info("%s: create sub task:%s", s.Info.Name, t.Name)
		r.Add(&t, s)
	}
	return nil
}

func (s *BatchTask) Stop() error {
	return nil
}

func NewBatchTask(handler common.TaskHandler, info *msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.BatchTask{}
	json.Unmarshal(data, &conf)

	return &BatchTask{
		TaskBase: common.NewTaskBase(info),
		config:   conf,
		handler:  handler,
	}
}
