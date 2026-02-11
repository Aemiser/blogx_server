package enum

type RoleType int8

const (
	UserRole    RoleType = iota + 1 // 用户角色
	AdminRole                       // 管理员角色
	visitorRole                     // 游客角色
)
