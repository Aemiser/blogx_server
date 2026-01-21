package log_service

import (
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"fmt"

	"github.com/gin-gonic/gin"
)

func NewLoginSuccess(c *gin.Context, loginType enum.LoginType) {
	ip := c.ClientIP()
	addr := core.GetIpAddr(ip)

	token := c.GetHeader("token")
	fmt.Printf("token: %s\n", token)
	//TODO :通过jwt获取用户ID,和用户信息
	userID := uint(1)
	username := "test"
	global.Db.Create(&models.LogModel{
		LogType:     enum.LoginLogType, //日志类型
		Title:       "登录成功",
		Content:     "",
		UserID:      userID,
		IP:          ip,
		Addr:        addr,
		LoginStatus: true,
		Username:    username,
		Pwd:         "-",
		LoginType:   loginType,
	})

}

func NewLoginFail(c *gin.Context, loginType enum.LoginType, msg, username, pwd string) {
	ip := c.ClientIP()
	addr := core.GetIpAddr(ip)

	global.Db.Create(&models.LogModel{
		LogType:     enum.LoginLogType, //日志类型
		Title:       "用户登录失败",
		Content:     msg,
		IP:          ip,
		Addr:        addr,
		LoginStatus: false,
		Username:    username,
		Pwd:         pwd,
		LoginType:   loginType,
	})

}
