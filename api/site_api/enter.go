package site_api

import (
	"blogx_server/common/res"

	"github.com/gin-gonic/gin"
)

type SiteApi struct {
}

func (SiteApi) SiteInfoView(c *gin.Context) {
	res.SuccessWithData("XX", c)
	return
}

type SiteUpdateRequest struct {
	Name string `json:"name"`
}

func (SiteApi) SiteUpdateView(c *gin.Context) {
	var req SiteUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithMsg(err.Error(), c)
	}
	res.SuccessWithMsg("更新成功", c)
	return
}
