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

	log.SetTitle("更新站点信息")
	var req SiteUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf(err.Error())
	}

	log.SetItemInfo("结构体", req)
	log.SetItemInfo("切片", []string{"a", "b"})
	log.SetItemInfo("Map", map[string]any{"a": "1", "b": "2"})
	log.SetItemInfo("字符串", "你好")
	log.SetItemInfo("数字", 123)

	c.JSON(200, gin.H{

		"code": 200,
		"data": gin.H{
			"name":    "site",
			"version": "1.0.0",
		},
	})
	return
}
