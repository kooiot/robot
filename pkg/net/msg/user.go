package msg

type Message struct {
	CTX  interface{} `json:"ctx"`
	Data []byte      `json:"data"`
}

type Login struct {
	ClientID string `json:"client_id"`
	User     string `json:"user"`
	Passwd   string `json:"passwd"`
	Hostname string `json:"hostname"`
	Hardware string `json:"hardware"`
	System   string `json:"system"`
}

type LoginResp struct {
	ClientID string `json:"client_id"`
	ID       int32  `json:"id"`
	Reason   string `json:"reason"`
}

type HeartBeat struct {
	ID   int32 `json:"id"`
	Time int64 `json:"time"`
}

type Logout struct {
	ClientID string `json:"client_id"`
	ID       int32  `json:"id"`
	Reason   string `json:"reason"`
}

type Response struct {
	Content string `json:"content"`
}
