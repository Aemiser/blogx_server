package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/pwd"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ResetPassowrdRequest struct {
	Pwd string `json:"pwd" binding:"required" `
}

func (UserApi) ResetPassowrdView(c *gin.Context) {
	var req ResetPassowrdRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	if !global.Config.Site.Login.EmailLogin {
		res.FailWithMsg("站点未启用邮箱注册", c)
		return
	}

	_email, _ := c.Get("email")
	email := _email.(string)

	var model models.UserModel
	err = global.Db.Take(&model, "email = ?", email).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		logrus.Errorf("用户不存在: %s", email)
		return
	}

	if model.RegisterSource != enum.RegisterSourceTypeEmail {
		res.FailWithMsg("非邮箱用户，不能重置密码", c)
		return
	}

	hashpwd, _ := pwd.GenerateHashPassword(req.Pwd)

	global.Db.Model(&model).Update("password", hashpwd)
	res.SuccessWithMsg("重置密码成功", c)

}
