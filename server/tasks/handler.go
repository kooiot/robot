package tasks

import (
	"encoding/json"
	"fmt"
	"path"
	"regexp"

	"github.com/Allenxuxu/gev"
	"github.com/kooiot/robot/pkg/net/msg"
	"github.com/kooiot/robot/pkg/net/protocol"
	"github.com/kooiot/robot/pkg/util/log"
	"github.com/kooiot/robot/server/common"
	"github.com/kooiot/robot/server/config"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type TaskInfo struct {
	Matches []config.AutoMatch
	Task    msg.BatchTask
}

type TaskHandler struct {
	base_path string
	autos     *config.AutoTasks
	tasks     []TaskInfo
}

func toStringMap(value interface{}) interface{} {
	switch value.(type) {
	case map[interface{}]interface{}:
		val := value.(map[interface{}]interface{})
		for k, v := range val {
			val[k] = toStringMap(v)

		}
		value = cast.ToStringMap(val)
	case map[string]interface{}:
		val := value.(map[string]interface{})
		for k, v := range val {
			val[k] = toStringMap(v)

		}
		value = val
	}

	return value
}

func parseTask(file_path string) (msg.BatchTask, error) {
	bt := msg.BatchTask{}

	v := viper.New()
	v.SetConfigFile(file_path)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	v.WatchConfig()

	if err := v.Unmarshal(&bt); err != nil {
		fmt.Println(err)
		return bt, err
	}
	new_bt := msg.BatchTask{}
	for _, v := range bt.Tasks {
		v.Option = toStringMap(v.Option)
		new_bt.Tasks = append(new_bt.Tasks, v)
	}
	log.Info("Task loaded: %#v", new_bt)
	return new_bt, nil
}

func (h *TaskHandler) Init(server common.Server) error {
	config_dir := server.ConfigDir()
	base_path := path.Join(config_dir, h.autos.Folder)
	h.base_path = base_path
	log.Info("Task loading from: %s", h.autos.Folder)

	for _, t := range h.autos.Autos {
		config_path := path.Join(base_path, t.Config)
		task, err := parseTask(config_path)
		if err != nil {
			log.Error("Task loading failed: %s", err)
		} else {
			h.tasks = append(h.tasks, TaskInfo{
				Task:    task,
				Matches: t.Matches,
			})
		}
	}

	return nil
}

func matchString(pat string, value string) bool {
	var m = regexp.MustCompile(pat)

	m_list := m.FindStringSubmatch(value)
	if m_list == nil {
		log.Info("Not matched %s - %s", pat, value)
	} else {
		log.Info("Task matched %s - %s", pat, value)
	}
	return m_list != nil
}

func (h *TaskHandler) AfterLogin(conn *gev.Connection, client *common.Client) {
	for _, t := range h.tasks {
		found := true
		for _, m := range t.Matches {
			switch m.Key {
			case "client_id":
				found = matchString(m.Match, client.Info.ClientID)
			case "hardware":
				found = matchString(m.Match, client.Info.Hardware)
			case "hostname":
				found = matchString(m.Match, client.Info.Hostname)
			case "user":
				found = matchString(m.Match, client.Info.User)
			default:
				found = false
			}

			if !found {
				break
			}
		}
		if found {
			log.Info("Fire task to: %s - %#v", client.Info.ClientID, t.Task)
			task := msg.Task{
				UUID:        uuid.NewV4().String(),
				Name:        "batch",
				Description: "Auto batch task",
				Option:      t.Task,
			}
			data, err := json.Marshal(&task)
			if err != nil {
				log.Error("failed encode resp: %s", err)
				log.Error("resp: %#v", task)
			} else {
				log.Info("Fire task json %s", data)
				conn.Send(protocol.PackMessage("task", data))
			}
		}
	}
}

func (h *TaskHandler) BeforeLogout(conn *gev.Connection, client *common.Client) {

}

func NewTaskHandler(autos *config.AutoTasks) *TaskHandler {
	return &TaskHandler{autos: autos}
}
