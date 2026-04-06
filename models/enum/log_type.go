package enum

type LogType int8

const (
	LoginLogType     = 1 // 登录日志
	ActionLogType    = 2 // 操作日志(默认)
	RuntimeLogType   = 3 // 运行日志
	OperationLogType = 4 // 管理员操作日志(写操作)
	QueryLogType     = 5 // 查询日志(GET请求)
)
