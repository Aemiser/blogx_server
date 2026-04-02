package user_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/user_service"
	"blogx_server/utils"
	"blogx_server/utils/email_store"
	"blogx_server/utils/pwd"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RegisterEmailRequest struct {
	EmailID string `json:"emailID" binding:"required" `
	Code    string `json:"code" binding:"required" `
	Pwd     string `json:"pwd" binding:"required" `
}

func (UserApi) RegisterEmail(c *gin.Context) {
	var req RegisterEmailRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	if !global.Config.Site.Login.EmailLogin {
		res.FailWithMsg("站点未启用邮箱注册", c)
		return
	}

	// 验证邮箱验证码并获取邮箱地址
	info, ok := email_store.Verify(req.EmailID, req.Code)
	if !ok {
		res.FailWithMsg("邮箱验证码错误", c)
		return
	}
	email := info.Email

	uname := fmt.Sprintf("b_%s", utils.GetRandomInDigital(4))
	unickname := fmt.Sprintf("邮箱用户%s", uname)
	pwd, err := pwd.GenerateHashPassword(req.Pwd)
	if err != nil {
		logrus.Errorf("生成密码失败: %s", err)
		res.FailWithError(err, c)
		return
	}

	//创建用户
	var newUser = models.UserModel{
		Model:          models.Model{},
		Username:       uname,
		Nickname:       unickname,
		RegisterSource: enum.RegisterSourceTypeEmail,
		Password:       pwd,
		Email:          email,
		Role:           enum.UserRole,
	}

	err = global.Db.Create(&newUser).Error
	if err != nil {
		res.FailWithMsg("邮箱注册失败", c)
		logrus.Errorf("创建用户失败: %s", err)
		return
	}

	// 颁发token
	token, err := jwts.GetToken(jwts.Claims{
		UserID:   newUser.ID,
		UserName: newUser.Username,
		Role:     newUser.Role,
	})
	if err != nil {
		res.FailWithMsg("邮箱登入失败", c)
		return
	}
	// 记入登入日志
	user_service.NewUserService(&newUser).UserLogin(c)
	res.SuccessWithData(token, c)
}
