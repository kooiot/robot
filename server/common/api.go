package common

import (
	"github.com/kooiot/robot/pkg/net/msg"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRole struct {
	RoleName string `json:"roleName"`
	Value    string `json:"value"`
}

type LoginResp struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	UUID     string `json:"uuid"`
	Token    string `json:"token"`
}

type UserInfo struct {
	UserID   uint32     `json:"userId"`
	Username string     `json:"username"`
	RealName string     `json:"realName"`
	Avatar   string     `json:"avatar"`
	Desc     string     `json:"desc"`
	HomePath string     `json:"homePath"`
	Roles    []UserRole `json:"roles"`
}

type MenuMeta struct {
	Title              string `json:"title"`
	Icon               string `json:"icon"`
	FrameSrc           string `json:"fromSrc"`
	CurrentActiveMenu  string `json:"currentActiveMenu"`
	HideMenu           bool   `json:"hideMenu"`
	HideChildrenInMenu bool   `json:"hideChildrenInMenu"`
	HideBreadcrumb     bool   `json:"hideBreadcrumb"`
	IgnoreKeepAlive    bool   `json:"ignoreKeepAlive"`
}

type Menu struct {
	Path      string   `json:"path"`
	Name      string   `json:"name"`
	Component string   `json:"component"`
	Redirect  string   `json:"redirect"`
	Meta      MenuMeta `json:"meta"`
	Children  []Menu   `json:"children"`
}

type RunResult struct {
	Hour    string `json:"hour"`
	Success uint32 `json:"success"`
	Fail    uint32 `json:"fail"`
}

type OnlineData struct {
	Hour  string `json:"hour"`
	Count uint32 `json:"count"`
}

type ErrorData struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Count uint32 `json:"count"`
}
type ErrorStatsData struct {
	Name  string `json:"name"`
	Count uint32 `json:"count"`
}

type DetailQuery struct {
	ID int32 `json:"id"`
}

type InfoQuery struct {
	ClientID string `json:"client_id"`
}

type ClientData struct {
	ID     int32     `json:"id"`
	Info   msg.Login `json:"info"`
	Online string    `json:"online"`
	Status string    `json:"status"`
}

type TestStatus struct {
	Today RunResult `json:"today"`
	Total RunResult `json:"total"`
}

type ServerStatus struct {
	Online  uint32 `json:"online"`
	Runing  uint32 `json:"running"`
	Total   uint32 `json:"total"`
	Success uint32 `json:"success"`
	Failed  uint32 `json:"failed"`
}

type RobotStats struct {
	UpdateTime   string           `json:"update_time"`
	Clients      []ClientData     `json:"clients"`
	RunToday     []RunResult      `json:"run_today"`
	ClientActive []RunResult      `json:"client_active"`
	ErrorTop     []ErrorData      `json:"error_top"`
	ErrorStats   []ErrorStatsData `json:"error_stats"`
	ClientLevel  []RunResult      `json:"client_level"`
	TestStatus   TestStatus       `json:"test_status"`
	ServerStatus ServerStatus     `json:"server_status"`
}
