package user_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/log_service"
	"blogx_server/service/user_service"
	"blogx_server/utils/pwd"
	"fmt"

	"github.com/gin-gonic/gin"
)

type PwdLoginRequest struct {
	Val      string `json:"val" binding:"required"`
	Password string `json:"pwd" binding:"required" `
}

func (UserApi) PwdLoginApi(c *gin.Context) {
	log := log_service.GetLog(c)
	log.SetLogType(enum.LoginLogType)
	log.SetTitle("<span style='color: #52c41a'>🔐 用户名密码登录</span>")
	log.ShowRequest()
	log.SetItem("登录方式", "用户名/邮箱 + 密码")

	req := middlerware.GetBind[PwdLoginRequest](c)
	log.SetItem("登录账号", req.Val)

	if !global.Config.Site.Login.UsernamePwdLogin {
		res.FailWithMsg("未启用用户名密码登录", c)
		return
	}
	// 查库
	var user models.UserModel
	err := global.Db.Take(&user, "(username = ? or email = ?) and password <> ''", req.Val, req.Val).Error
	if err != nil {
		log.SetItem("登录结果", "<span style='color: #ff4d4f'>❌ 用户不存在</span>")
		res.FailWithMsg("用户名或密码错误", c)
		return
	}

	//验证密码
	if !pwd.CompareHashAndPassword(user.Password, req.Password) {
		fmt.Println(user.Password)
		fmt.Println(req.Password)
		log.SetItem("登录结果", "<span style='color: #ff4d4f'>❌ 密码错误</span>")
		res.FailWithMsg("用户名或密码错误", c)
		return
	}

	log.SetItem("用户ID", fmt.Sprintf("<span style='color: #1890ff'>%d</span>", user.ID))
	log.SetItem("用户昵称", user.Nickname)
	log.SetItem("登录结果", "<span style='color: #52c41a'>✅ 登录成功</span>")

	// 颁发token
	token, _ := jwts.GetToken(jwts.Claims{
		UserID:   user.ID,
		UserName: user.Username,
		Role:     user.Role,
	})

	user_service.NewUserService(&user).UserLogin(c)
	res.SuccessWithData(token, c)
}
