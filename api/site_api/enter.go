package site_api

import (
	"blogx_server/models/enum"
	"blogx_server/service/log_service"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SiteApi struct {
}

func (SiteApi) SiteInfoView(c *gin.Context) {
	fmt.Printf("1")
	log_service.NewLoginSuccess(c, enum.UserPwdLoginType)
	log_service.NewLoginFail(c, enum.UserPwdLoginType, "用户名不存在", "test", "1234")
	c.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"name":    "site",
			"version": "1.0.0",
		},
	})
	return
}

type SiteUpdateRequest struct {
	Name string `json:"name"`
}

func (SiteApi) SiteUpdateView(c *gin.Context) {
	log := log_service.GetLog(c)

	log.ShowRequest()
	log.ShowResponse()

	var req SiteUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf(err.Error())
	}
	c.JSON(200, gin.H{

		"code": 200,
		"data": gin.H{
			"name":    "site",
			"version": "1.0.0",
		},
	})
	return
}
