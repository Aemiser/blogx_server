package log_service

import (
	"blogx_server/common/jwts"
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewLoginSuccess(c *gin.Context, loginType enum.LoginType) {
	ip := c.ClientIP()
	addr := core.GetIpAddr(ip)

	claim, err := jwts.ParseTokenByGin(c)
	userID := uint(0)
	username := ""
	if err == nil && claim != nil {
		username = claim.Claims.UserName
		userID = claim.Claims.UserID
	}
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
	logrus.Info("日志创建成功")
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
