package user_service

import (
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (u UserService) UserLogin(c *gin.Context) {
	ip := c.ClientIP()
	addr := core.GetIpAddr(ip)
	ua := c.GetHeader("User-Agent")
	err := global.Db.Create(&models.UserLoginModel{
		UserID: u.ID,
		IP:     ip,
		Addr:   addr,
		UA:     ua,
	}).Error

	if err != nil {
		logrus.Errorf("用户登录记录失败: %s", err)
	}

}
