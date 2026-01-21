package enum

type LoginType int8

const (
	UserPwdLoginType = 1
	QQLoginType      = 2
	EmailLoginType   = 3
	PhoneLoginType   = 4
	WechatLoginType  = 5
)
