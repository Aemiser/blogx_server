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
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required" label:"年龄"`
}

func (SiteApi) SiteUpdateView(c *gin.Context) {
	var req SiteUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	res.SuccessWithMsg("更新成功", c)
	return
}
