package tasks

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"io"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/shutdown"
	"github.com/kooiot/robot/pkg/util/xlog"
)

type TaskStoreMsgType int32

const (
	MSG_OPEN   TaskStoreMsgType = 0
	MSG_CLOSE  TaskStoreMsgType = 1
	MSG_UPDATE TaskStoreMsgType = 2
	MSG_RESULT TaskStoreMsgType = 4
)

type TaskStoreMsg struct {
	Type     TaskStoreMsgType
	ClientID string
	Msg      interface{}
}

type TaskStoreInfo struct {
	ClientID string
	Output   string
	Tasks    []msg.TaskResult
}

type TaskStore struct {
	onTaskInsert func(client_id string, task *msg.Task)
	onTaskUpdate func(client_id string, task *msg.Task)
	onTaskFinish func(client_id string, result *msg.TaskResult)
	onOpen       func(client_id string)
	onClose      func(client_id string, err error)

	ctx        context.Context
	xl         *xlog.Logger
	output_dir string

	clients map[string]*TaskStoreInfo

	msg_chn         chan TaskStoreMsg
	worker_shutdown *shutdown.Shutdown
}

func (t *TaskStore) OnTaskInsert(f func(client_id string, task *msg.Task)) {
	t.onTaskInsert = f
}

func (t *TaskStore) OnTaskUpdate(f func(client_id string, task *msg.Task)) {
	t.onTaskUpdate = f
}

func (t *TaskStore) OnTaskFinish(f func(client_id string, result *msg.TaskResult)) {
	t.onTaskFinish = f
}

func (t *TaskStore) OnOpen(f func(client_id string)) {
	t.onOpen = f
}

func (t *TaskStore) OnClose(f func(client_id string, err error)) {
	t.onClose = f
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

func (t *TaskStore) dump_tasks(info *TaskStoreInfo) {
	xl := xlog.FromContextSafe(t.ctx)
	data, err := json.Marshal(info.Tasks)
	if err != nil {
		xl.Error("JSON.Marshal failure: %s", err.Error())
	} else {
		err = os.WriteFile(info.Output, data, 0644)
		if err != nil {
			xl.Error("os.WriteFile to: %s failure: %s", info.Output, err.Error())
		}
	}
}

func (t *TaskStore) worker() {
	xl := t.xl
	defer func() {
		if err := recover(); err != nil {
			xl.Error("panic error: %v", err)
			xl.Error(string(debug.Stack()))
		}
	}()
	defer t.worker_shutdown.Done()

	xl.Debug("worker start.......")
	for m := range t.msg_chn {
		xl.Debug("Process.......")
		switch m.Type {
		case MSG_OPEN:
			t._Open(m.ClientID)
		case MSG_CLOSE:
			t._Close(m.ClientID, m.Msg.(error))
		case MSG_UPDATE:
			t._TaskUpdate(m.ClientID, m.Msg.(*msg.Task))
		case MSG_RESULT:
			t._TaskResult(m.ClientID, m.Msg.(*msg.TaskResult))
		}
	}
	xl.Debug("worker quit.......")
}

func (t *TaskStore) Start() {
	go t.worker()
}

func (t *TaskStore) Stop() {
	close(t.msg_chn)

	t.worker_shutdown.WaitDone()
}

func (t *TaskStore) Open(client_id string) {
	t.msg_chn <- TaskStoreMsg{
		Type:     MSG_OPEN,
		ClientID: client_id,
	}
}

func (t *TaskStore) _Open(client_id string) {
	xl := t.xl
	now := strconv.FormatInt(time.Now().Unix(), 10)
	output_path := path.Join(t.output_dir, client_id+"_"+now+".json")
	xl.Info("Task result save to: %s", output_path)
	t.clients[client_id] = &TaskStoreInfo{
		ClientID: client_id,
		Output:   output_path,
	}
	t.onOpen(client_id)
}

func (t *TaskStore) Close(client_id string, err error) {
	t.msg_chn <- TaskStoreMsg{
		Type:     MSG_CLOSE,
		ClientID: client_id,
		Msg:      err,
	}
}

func (t *TaskStore) _Close(client_id string, err error) {
	xl := t.xl
	info, ok := t.clients[client_id]
	if !ok {
		xl.Error("task for client:%d missing", client_id)
	}
	defer delete(t.clients, client_id)

	t.dump_tasks(info)
	t.onClose(client_id, err)
}

func (t *TaskStore) TaskUpdate(client_id string, task *msg.Task) {
	t.msg_chn <- TaskStoreMsg{
		Type:     MSG_UPDATE,
		ClientID: client_id,
		Msg:      task,
	}
}

func (t *TaskStore) _TaskUpdate(client_id string, task *msg.Task) {
	xl := t.xl
	info, ok := t.clients[client_id]
	if !ok {
		xl.Error("task for client:%d missing", client_id)
	}
	for i, v := range info.Tasks {
		if v.Task.UUID == task.UUID {
			info.Tasks[i].Task = *task
			t.dump_tasks(info)
			t.onTaskUpdate(client_id, task)
			return
		}
	}
	// Insert New
	info.Tasks = append(info.Tasks, msg.TaskResult{
		Task: *task,
	})

	t.dump_tasks(info)
	t.onTaskInsert(client_id, task)
}

func (t *TaskStore) TaskResult(client_id string, result *msg.TaskResult) {
	t.msg_chn <- TaskStoreMsg{
		Type:     MSG_RESULT,
		ClientID: client_id,
		Msg:      result,
	}
}

func (t *TaskStore) _TaskResult(client_id string, result *msg.TaskResult) {
	xl := t.xl
	info, ok := t.clients[client_id]
	if !ok {
		xl.Error("task for client:%d missing", client_id)
	}
	for i, v := range info.Tasks {
		if v.Task.UUID == result.Task.UUID {
			info.Tasks[i] = *result
			t.dump_tasks(info)
			t.onTaskUpdate(client_id, &result.Task)
			t.onTaskFinish(client_id, result)
			return
		}
	}
	t.xl.Error("task not found for result: %#v", result)
}

func NewTaskStore(ctx context.Context, output string) *TaskStore {
	store := &TaskStore{
		ctx:             ctx,
		xl:              xlog.FromContextSafe(ctx),
		msg_chn:         make(chan TaskStoreMsg, 100),
		clients:         make(map[string]*TaskStoreInfo),
		worker_shutdown: shutdown.New(),
		output_dir:      output,
	}

	store.OnOpen(func(client_id string) {})
	store.OnClose(func(client_id string, err error) {})
	store.OnTaskInsert(func(client_id string, task *msg.Task) {})
	store.OnTaskUpdate(func(client_id string, task *msg.Task) {})
	store.OnTaskFinish(func(client_id string, result *msg.TaskResult) {})

	return store
}
