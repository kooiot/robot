package msg

type Task struct {
	UUID        string      `json:"uuid"`
	Name        string      `json:"name"`
	Description string      `json:"desc"`
	Option      interface{} `json:"option"`
}

type BatchTask struct {
	Tasks []Task `json:"tasks"`
}

type TaskResult struct {
	UUID   string `json:"uuid"`
	Result bool   `json:"result"`
	Info   string `json:"info"`
}

type SerialTask struct {
	SrcPort    string `json:"src"`
	DestPort   string `json:"dst"`
	Baudrate   int    `json:"baudrate"`
	Count      int    `json:"count"`
	MaxMsgSize int    `json:"max_msg_size"`
}

type RTCTask struct {
	File string `json:"file"`
}

type USBTask struct {
	IDS   []string `json:"ids"`
	Reset string   `json:"reset"`
	Power string   `json:"power"`
}

type ModemTask struct {
	PingAddr string  `json:"ping_addr"`
	USB      USBTask `json:"usb"`
}

type NamedGPIOTask struct {
	Name string `json:"name"`
}
