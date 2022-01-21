package msg

type StatusType int32

const (
	ST_NEW    StatusType = 0
	ST_RUN    StatusType = 1
	ST_FAILED StatusType = 2
	ST_DONE   StatusType = 4
	ST_SPAWN  StatusType = 8
)

type Task struct {
	UUID        string      `mapstructure:"uuid" json:"uuid"`
	ParentUUID  string      `mapstructure:"parent_uuid" json:"parent_uuid"`
	Status      StatusType  `mapstructure:"status" json:"status"`
	ID          string      `mapstructure:"id" json:"id"`
	Task        string      `mapstructure:"task" json:"task"`
	Description string      `mapstructure:"desc" json:"desc"`
	Option      interface{} `mapstructure:"option" json:"option"`
	Depends     []string    `mapstructure:"depends" json:"depends"`
}

type TaskResultDetail struct {
	Result bool        `mapstructure:"result" json:"result"`
	Info   string      `mapstructure:"info" json:"info"`
	Detail interface{} `mapstructure:"detail" json:"detail"`
}

type TaskResult struct {
	Task   Task             `mapstructure:"task" json:"task"`
	Detail TaskResultDetail `mapstructure:"detail" json:"detail"`
}

type BatchTask struct {
	Tasks []Task `mapstructure:"tasks" json:"tasks"`
}

type SerialTask struct {
	SrcPort    string `mapstructure:"src" json:"src"`
	DstPort    string `mapstructure:"dst" json:"dst"`
	Baudrate   int    `mapstructure:"baudrate" json:"baudrate"`
	Count      int    `mapstructure:"count" json:"count"`
	MaxMsgSize int    `mapstructure:"max_msg_size" json:"max_msg_size"`
}

type RTCTask struct {
	File string `mapstructure:"file" json:"file"`
}

type USBETHTask struct {
	Name  string `mapstructure:"name" json:"name"`
	Reset string `mapstructure:"reset" json:"reset"`
	Power string `mapstructure:"power" json:"power"`
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

type EthernetTask struct {
	PingAddr string   `mapstructure:"ping_addr" json:"ping_addr"`
	Init     []string `mapstructure:"init" json:"init"`
}

type NamedGPIOTask struct {
	Name string `mapstructure:"name" json:"name"`
}

type LedsTask struct {
	Leds  []string `mapstructure:"leds" json:"leds"`
	Count int      `mapstructure:"count" json:"count"`
	Span  int      `mapstructure:"span" json:"span"`
}

type DoneTask struct {
	Leds []string `mapstructure:"leds" json:"leds"`
	Halt bool     `mapstructure:"halt" json:"halt"`
}
