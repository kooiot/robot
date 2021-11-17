package common

type StreamResult struct {
	Count     int     `json:"count"`
	Failed    int     `json:"failed"`
	Passed    int     `json:"passed"`
	SendBytes int     `json:"send_bytes"`
	RecvBytes int     `json:"recv_bytes"`
	ErrBytes  int     `json:"err_bytes"`
	SendSpeed float64 `json:"send_speed"`
	RecvSpeed float64 `json:"recv_speed"`
}
