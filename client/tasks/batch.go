package tasks

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kooiot/robot/client/common"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type BatchTask struct {
	common.TaskBase
	ctx     context.Context
	config  msg.BatchTask
	handler common.TaskHandler
}

func init() {
	RegisterTask("batch", NewBatchTask)
}

func (s *BatchTask) Run() (interface{}, error) {
	xl := xlog.FromContextSafe(s.ctx)

	r, ok := s.handler.(*Runner)
	if !ok {
		return nil, errors.New("error object")
	}
	for _, t := range s.config.Tasks {
		xl.Debug("%s: create sub task:%s", s.Info.Task, t.Task)
		r.Add(t, s)
	}
	return "Done", nil
}

func NewBatchTask(ctx context.Context, handler common.TaskHandler, info msg.Task, parent common.Task) common.Task {
	data, _ := json.Marshal(info.Option)

	conf := msg.BatchTask{}
	json.Unmarshal(data, &conf)

	xl := xlog.FromContextSafe(ctx).Spawn().AppendPrefix("Task.Batch")
	return &BatchTask{
		TaskBase: common.NewTaskBase(info),
		ctx:      xlog.NewContext(ctx, xl),
		config:   conf,
		handler:  handler,
	}
}
