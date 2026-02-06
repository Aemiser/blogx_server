package user_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models/enum"
	"blogx_server/utils/pwd"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UpdatePasswordRequest struct {
	OldPwd string `json:"oldPwd" binding:"required"`
	Pwd    string `json:"pwd" binding:"required"`
}

func (UserApi) UpdatePasswordView(c *gin.Context) {
	var req UpdatePasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	// 获取当前用户
	claims := jwts.GetClaimsByGin(c)
	user, err := claims.GetUser()
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	if user.RegisterSource != enum.RegisterSourceTypeEmail || user.Email == "" {
		res.FailWithMsg("仅支持邮箱和绑定了邮箱的用户修改密码", c)
		return
	}

	// 校验之前的密码
	if !pwd.CompareHashAndPassword(user.Password, req.OldPwd) {
		res.FailWithMsg("旧密码错误", c)
		return
	}

	//获取密码哈希
	_hash, _ := pwd.GenerateHashPassword(req.Pwd)

	//修改密码
	err = global.Db.Model(&user).Update("password", _hash).Error
	if err != nil {
		logrus.Errorf("修改密码失败: %v", err)
		res.FailWithMsg("修改密码失败", c)
		return
	}
	res.SuccessWithMsg("修改密码成功", c)
}
