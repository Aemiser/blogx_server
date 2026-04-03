package res

type Code int

const (
	SuccessCode     Code = 0
	FailValueCode   Code = 1001
	FailServiceCode Code = 1002
	SysFail         Code = 1003
	SysUnauthorized Code = 1004
	SysForbidden    Code = 1005
	SysNotFound     Code = 1006
)

var CodeMap = map[Code]string{}

func RegisterCode(c Code, msg string) {
	CodeMap[c] = msg
}

func GetMsg(code Code) string {
	if msg, ok := CodeMap[code]; ok {
		return msg
	}
	return "未知错误"
}

func (c Code) String() string {
	if msg, ok := CodeMap[c]; ok {
		return msg
	}
	return "未知错误"
}

func InitSysCode() {
	RegisterCode(SuccessCode, "成功")
	RegisterCode(FailValueCode, "参数错误")
	RegisterCode(FailServiceCode, "服务错误")
	RegisterCode(SysFail, "操作失败")
	RegisterCode(SysUnauthorized, "未授权")
	RegisterCode(SysForbidden, "禁止访问")
	RegisterCode(SysNotFound, "资源不存在")
}
