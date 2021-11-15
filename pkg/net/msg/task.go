package msg

type Task struct {
	UUID    string      `json:"uuid"`
	Name    string      `json:"name"`
	Command string      `json:"command"`
	Params  interface{} `json:"params"`
}

type BatchTask struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

type TaskResult struct {
	UUID   string `json:"uuid"`
	Result bool   `json:"result"`
	Info   string `json:"info"`
}
