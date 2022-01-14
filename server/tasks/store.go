package tasks

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type TaskStore struct {
	onTaskInsert func(task *msg.Task)
	onTaskUpdate func(task *msg.Task)
	onTaskFinish func(task *msg.Task, result *msg.TaskResult)
	onOpen       func()
	onClose      func(err error)

	ctx         context.Context
	xl          *xlog.Logger
	output_dir  string
	output_path string
	ClientID    string
	Tasks       []msg.Task
}

func (t *TaskStore) OnTaskInsert(f func(task *msg.Task)) {
	t.onTaskInsert = f
}

func (t *TaskStore) OnTaskUpdate(f func(task *msg.Task)) {
	t.onTaskUpdate = f
}

func (t *TaskStore) OnTaskFinish(f func(task *msg.Task, result *msg.TaskResult)) {
	t.onTaskFinish = f
}

func (t *TaskStore) OnOpen(f func()) {
	t.onOpen = f
}

func (t *TaskStore) OnClose(f func(error)) {
	t.onClose = f
}

func (t *TaskStore) Open() {
	xl := xlog.FromContextSafe(t.ctx)
	now := strconv.FormatInt(time.Now().Unix(), 10)
	t.output_path = path.Join(t.output_dir, t.ClientID+"_"+now+".json")
	xl.Info("Task result save to: %s", t.output_path)

	t.onOpen()
}

func DoZlibCompress(dataSrc []byte) []byte {
	var out bytes.Buffer
	w := zlib.NewWriter(&out)
	w.Write(dataSrc)
	w.Close()
	return out.Bytes()
}

func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

func (t *TaskStore) dump_tasks() {
	xl := xlog.FromContextSafe(t.ctx)
	data, err := json.Marshal(t.Tasks)
	if err != nil {
		xl.Error("JSON.Marshal failure: %s", err.Error())
	} else {
		xl.Debug("Task result save to: %s", t.output_path)
		//err = os.WriteFile(t.output_path, DoZlibCompress(data), 0644)
		err = os.WriteFile(t.output_path, data, 0644)
		if err != nil {
			xl.Error("os.WriteFile to: %s failure: %s", t.output_path, err.Error())
		}
	}
}

func (t *TaskStore) Close(err error) {
	t.dump_tasks()
	t.onClose(err)
}

func (t *TaskStore) TaskUpdate(task *msg.Task) {
	for i, v := range t.Tasks {
		if v.UUID == task.UUID {
			t.Tasks[i] = *task
			t.onTaskUpdate(task)

			t.dump_tasks()
			return
		}
	}
	// Insert New
	t.Tasks = append(t.Tasks, *task)
	t.onTaskInsert(task)

	t.dump_tasks()
}

func (t *TaskStore) TaskResult(result *msg.TaskResult) {
	for _, v := range t.Tasks {
		if v.UUID == result.Task.UUID {
			t.dump_tasks()
			t.onTaskFinish(&v, result)
			return
		}
	}
	t.xl.Error("task not found for result: %#v", result)
}

func NewTaskStore(ctx context.Context, client_id string, output string) *TaskStore {
	store := &TaskStore{
		ctx:        ctx,
		xl:         xlog.FromContextSafe(ctx),
		ClientID:   client_id,
		output_dir: output,
	}

	store.OnOpen(func() {})
	store.OnClose(func(err error) {})
	store.OnTaskInsert(func(task *msg.Task) {})
	store.OnTaskUpdate(func(task *msg.Task) {})
	store.OnTaskFinish(func(task *msg.Task, result *msg.TaskResult) {})

	return store
}
