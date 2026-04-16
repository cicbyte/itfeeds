package init

import "github.com/gogf/gf/v2/frame/g"

// InitStatusReq 查询初始化状态
type InitStatusReq struct {
	g.Meta `path:"/init/status" method:"get" tags:"初始化" summary:"查询初始化状态"`
}

// InitStatusRes 初始化状态响应
type InitStatusRes struct {
	g.Meta `mime:"application/json"`
	Initialized bool `json:"initialized"`
}

// InitTestReq 测试数据库连接请求
type InitTestReq struct {
	g.Meta   `path:"/init/test-connection" method:"post" tags:"初始化" summary:"测试数据库连接"`
	Host     string `json:"host" v:"required#数据库主机不能为空"`
	Port     int    `json:"port" v:"required|min:1|max:65535#端口不能为空|端口范围1-65535"`
	User     string `json:"user" v:"required#用户名不能为空"`
	Password string `json:"password"`
	Database string `json:"database" v:"required#数据库名不能为空"`
}

// InitTestRes 测试连接响应
type InitTestRes struct {
	g.Meta  `mime:"application/json"`
	Success bool   `json:"success"`
	Version string `json:"version,omitempty"`
	Error   string `json:"error,omitempty"`
}

// InitSetupReq 执行初始化请求
type InitSetupReq struct {
	g.Meta   `path:"/init/setup" method:"post" tags:"初始化" summary:"执行系统初始化"`
	Host     string `json:"host" v:"required#数据库主机不能为空"`
	Port     int    `json:"port" v:"required|min:1|max:65535#端口不能为空|端口范围1-65535"`
	User     string `json:"user" v:"required#用户名不能为空"`
	Password string `json:"password"`
	Database string `json:"database" v:"required#数据库名不能为空"`
}

// InitSetupRes 执行初始化响应
type InitSetupRes struct {
	g.Meta  `mime:"application/json"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
