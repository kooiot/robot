package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/util/xlog"
	"github.com/kooiot/robot/server/common"
)

var match_topic = regexp.MustCompile(`^(.+)_(\d+)\.json$`)

type ResultStats struct {
	ctx            context.Context
	xl             *xlog.Logger
	output_dir     string
	max_retain_day uint32
	lock           sync.RWMutex
	stats          common.RobotStats
	clients        map[int32]TaskStoreInfo
}

func (t *ResultStats) load_output_files() {
	xl := xlog.FromContextSafe(t.ctx)

	err := filepath.Walk(t.output_dir,
		func(path string, f os.FileInfo, err error) error {
			if f == nil {
				xl.Error(err.Error())
				return err
			}
			if f.IsDir() {
				return nil
			}

			xl.Trace("loading file %s %s", path, filepath.Base(path))

			tm := match_topic.FindStringSubmatch(filepath.Base(path))
			if tm == nil {
				return errors.New("topic match failed")
			}
			client_id := tm[1]
			timestamp_str, _ := strconv.ParseInt(tm[2], 10, 64)
			timestamp := time.Unix(timestamp_str, 0)
			xl.Trace("Loading result: %s %s", client_id, timestamp.String())

			if time.Since(timestamp) > time.Hour*time.Duration(t.max_retain_day)*24 {
				// Remove test files
				xl.Info("Removing file %s", filepath.Base(path))
				os.Remove(path)
				return nil
			}

			data, err := ioutil.ReadFile(path)
			if err != nil {
				// Error loading json file
				xl.Error("Loading file %s error %s", filepath.Base(path), err.Error())
				return nil
			}

			info := TaskStoreInfo{
				ClientID: client_id,
				Output:   path,
			}
			if err := json.Unmarshal(data, &info.Tasks); err != nil {
				// Error loading json file
				xl.Error("Loading file %s error %s", filepath.Base(path), err.Error())
				return nil
			}
			t._push_result(&info, timestamp)

			// Only Loading
			return nil
		})
	if err != nil {
		xl.Error(err.Error())
	}
}

func (t *ResultStats) _push_result(info *TaskStoreInfo, timestamp time.Time) {
	// now := time.Now()
	// today_begin := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// is_today := timestamp.Sub(today_begin) > 0
	is_today := time.Since(timestamp) <= time.Hour*24

	t.lock.Lock()
	defer t.lock.Unlock()

	for _, task := range info.Tasks {
		// Skip the spawned task??
		if task.Task.ParentUUID != "" {
			continue
		}

		t.stats.UpdateTime = time.Now().Format("2006-01-02 15:04:05")

		if task.Task.Status == msg.ST_DONE {
			t.stats.TestStatus.Total.Success = t.stats.TestStatus.Total.Success + 1
			if is_today {
				t.stats.TestStatus.Today.Success = t.stats.TestStatus.Today.Success + 1
			}
		} else {
			t.stats.TestStatus.Total.Fail = t.stats.TestStatus.Total.Fail + 1
			if is_today {
				t.stats.TestStatus.Today.Fail = t.stats.TestStatus.Today.Fail + 1
			}
		}
	}
}

func (t *ResultStats) Init() {
	t.load_output_files()
}

func (t *ResultStats) onTaskResult(client_id string, result *msg.TaskResult) {
	info := TaskStoreInfo{
		ClientID: client_id,
	}
	info.Tasks = append(info.Tasks, *result)

	t._push_result(&info, time.Now())
}

func (t *ResultStats) onOpen(client_id string, info *common.ClientData) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.stats.Clients = append(t.stats.Clients, *info)
	t.stats.ServerStatus.Total = t.stats.ServerStatus.Total + 1
	t.stats.ServerStatus.Runing = t.stats.ServerStatus.Runing + 1
	t.stats.ServerStatus.Online = t.stats.ServerStatus.Online + 1
}
func (t *ResultStats) onClose(client_id string, err error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.stats.ServerStatus.Online = t.stats.ServerStatus.Online - 1
	t.stats.ServerStatus.Runing = t.stats.ServerStatus.Runing - 1
	for index, client := range t.stats.Clients {
		if client.Info.ClientID == client_id {
			t.stats.Clients = append(t.stats.Clients[:index], t.stats.Clients[index+1:]...)
		}
	}
	for index, client := range t.clients {
		if client.ClientID == client_id {
			delete(t.clients, index)
		}
	}
}

func (t *ResultStats) UpdateClient(client_id string, info *TaskStoreInfo) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.clients[info.ID] = *info
}

func (t *ResultStats) GetStats() common.RobotStats {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.stats
}

func (t *ResultStats) GetDetail(client_id int32) TaskStoreInfo {
	t.lock.RLock()
	defer t.lock.RUnlock()
	// xl := xlog.FromContextSafe(t.ctx)
	// xl.Trace("Get detail for client %d", client_id)
	// xl.Trace("Clients %#v", t.clients)

	return t.clients[client_id]
}

func (t *ResultStats) GetInfo(client_id string) common.ClientData {
	t.lock.RLock()
	defer t.lock.RUnlock()

	for _, client := range t.stats.Clients {
		if client.Info.ClientID == client_id {
			return client
		}
	}
	return common.ClientData{}
}

func NewResultStats(ctx context.Context, output string, max_retain_day uint32) *ResultStats {
	stats := &ResultStats{
		ctx:            ctx,
		xl:             xlog.FromContextSafe(ctx),
		max_retain_day: max_retain_day,
		output_dir:     output,
		stats:          common.RobotStats{},
		clients:        map[int32]TaskStoreInfo{},
	}

	return stats
}
