package res

const (
	FocusCannotSelf     Code = 8001
	FocusUserNotFound   Code = 8002
	FocusAlreadyExists  Code = 8003
	FocusRestoreFailed  Code = 8004
	FocusCannotSelfUn   Code = 8005
	UnFocusUserNotFound Code = 8006
	UnFocusNotFollowed  Code = 8007
	FocusUserConfigNot  Code = 8008
	FocusUserNotOpen    Code = 8009
	FocusLoginRequired  Code = 8010
	FansUserConfigNot   Code = 8011
	FansUserNotOpen     Code = 8012
)

func InitFocusCode() {
	RegisterCode(FocusCannotSelf, "不能关注自己")
	RegisterCode(FocusUserNotFound, "关注用户不存在")
	RegisterCode(FocusAlreadyExists, "已经关注了")
	RegisterCode(FocusRestoreFailed, "关注恢复失败")
	RegisterCode(FocusCannotSelfUn, "不能取关自己")
	RegisterCode(UnFocusUserNotFound, "取关用户不存在")
	RegisterCode(UnFocusNotFollowed, "未关注该用户")
	RegisterCode(FocusUserConfigNot, "用户配置信息不存在")
	RegisterCode(FocusUserNotOpen, "用户未开放我的关注")
	RegisterCode(FocusLoginRequired, "请登录")
	RegisterCode(FansUserConfigNot, "用户配置信息不存在")
	RegisterCode(FansUserNotOpen, "用户未开放我的粉丝")
}
