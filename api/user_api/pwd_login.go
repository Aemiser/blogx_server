package user_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/utils/pwd"
	"fmt"

	"github.com/gin-gonic/gin"
)

type PwdLoginRequest struct {
	Val      string `json:"val" binding:"required"`
	Password string `json:"pwd" binding:"required" `
}

func (UserApi) PwdLoginApi(c *gin.Context) {
	var req PwdLoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	if !global.Config.Site.Login.UsernamePwdLogin {
		res.FailWithMsg("未启用用户名密码登录", c)
		return
	}
	// 查库
	var user models.UserModel
	err = global.Db.Take(&user, "(username = ? or email = ?) and password <> ''", req.Val, req.Val).Error
	if err != nil {
		res.FailWithMsg("用户名或密码错误", c)
		return
	}

	//验证密码
	if !pwd.CompareHashAndPassword(user.Password, req.Password) {
		fmt.Println(user.Password)
		fmt.Println(req.Password)
		res.FailWithMsg("用户名或密码错误", c)
		return
	}

	// 颁发token
	token, _ := jwts.GetToken(jwts.Claims{
		UserID:   user.ID,
		UserName: user.Username,
		Role:     user.Role,
	})
	res.SuccessWithData(token, c)
}
