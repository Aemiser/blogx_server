package res

const (
	UserNotFound        Code = 4001
	UserAlreadyExists   Code = 4002
	UserDisabled        Code = 4003
	UserPasswordError   Code = 4004
	UserTokenExpired    Code = 4005
	UserTokenInvalid    Code = 4006
	UserNoPermission    Code = 4007
	UserNameDuplicate   Code = 4008
	UserNameModifyLimit Code = 4009
	UserQQNotAllowed    Code = 4010
	UserInfoUpdateFail  Code = 4011
	UserNotOpenCollect  Code = 4012
)

func InitUserCode() {
	RegisterCode(UserNotFound, "用户不存在")
	RegisterCode(UserAlreadyExists, "用户已存在")
	RegisterCode(UserDisabled, "用户已被禁用")
	RegisterCode(UserPasswordError, "密码错误")
	RegisterCode(UserTokenExpired, "登录已过期")
	RegisterCode(UserTokenInvalid, "无效的登录")
	RegisterCode(UserNoPermission, "无权限")
	RegisterCode(UserNotOpenCollect, "用户未开放收藏")
	RegisterCode(UserNameDuplicate, "该用户名已被使用")
	RegisterCode(UserNameModifyLimit, "用户名30天内只能修改一次")
	RegisterCode(UserQQNotAllowed, "QQ用户不允许修改昵称和头像")
	RegisterCode(UserInfoUpdateFail, "用户信息修改失败")
}
