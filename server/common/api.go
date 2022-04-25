package common

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
	Success uint32 `json:"success"`
	Fail    uint32 `json:"fail"`
}

type OnlineData struct {
	Hour  uint32 `json:"hour"`
	Count uint32 `json:"count"`
}

type ErrorData struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Count uint32 `json:"count"`
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
	UpdateTime   string       `json:"update_time"`
	RunToday     []RunResult  `json:"run_today"`
	ClientActive []RunResult  `json:"client_active"`
	OrderSource  []RunResult  `json:"order_source"`
	ClientLevel  []RunResult  `json:"client_level"`
	ErrorTop     []ErrorData  `json:"error_top"`
	TestStatus   TestStatus   `json:"test_status"`
	ServerStatus ServerStatus `json:"server_status"`
}
