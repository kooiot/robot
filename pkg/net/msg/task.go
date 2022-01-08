package msg

type Task struct {
	UUID        string      `mapstructure:"uuid" json:"uuid"`
	Name        string      `mapstructure:"name" json:"name"`
	Description string      `mapstructure:"desc" json:"desc"`
	Option      interface{} `mapstructure:"option" json:"option"`
}

type BatchTask struct {
	Tasks []Task `mapstructure:"tasks" json:"tasks"`
}

type TaskResult struct {
	UUID   string `mapstructure:"uuid" json:"uuid"`
	Result bool   `mapstructure:"result" json:"result"`
	Info   string `mapstructure:"info" json:"info"`
}

type SerialTask struct {
	SrcPort    string `mapstructure:"src" json:"src"`
	DestPort   string `mapstructure:"dst" json:"dst"`
	Baudrate   int    `mapstructure:"baudrate" json:"baudrate"`
	Count      int    `mapstructure:"count" json:"count"`
	MaxMsgSize int    `mapstructure:"max_msg_size" json:"max_msg_size"`
}

type RTCTask struct {
	File string `mapstructure:"file" json:"file"`
}

type USBTask struct {
	IDS   []string `mapstructure:"ids" json:"ids"`
	Reset string   `mapstructure:"reset" json:"reset"`
	Power string   `mapstructure:"power" json:"power"`
}

type ModemTask struct {
	PingAddr string  `mapstructure:"ping_addr" json:"ping_addr"`
	USB      USBTask `mapstructure:"usb" json:"usb"`
}

type NamedGPIOTask struct {
	Name string `mapstructure:"name" json:"name"`
}

type DoneTask struct {
	Leds []string `mapstructure:"leds" json:"leds"`
}
