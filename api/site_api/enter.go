package site_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"fmt"

	"github.com/gin-gonic/gin"
)

type SiteApi struct {
}

type SiteInfoRequest struct {
	Name string `json:"name" form:"name" uri:"name"`
}

func (SiteApi) SiteInfoView(c *gin.Context) {
	var req SiteInfoRequest
	if err := c.ShouldBindUri(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	fmt.Println("req:", req)
	var data any
	if req.Name == "site" {
		fmt.Println("site")
		data = global.Config.Site
		res.SuccessWithData(data, c)
		return
	}

	// 判断管理员
	middlerware.AdminMiddleware(c)
	_, ok := c.Get("claims")
	if !ok {
		return
	}

	switch req.Name {
	case "email":
		data = global.Config.Email
	case "qq":
		data = global.Config.QQ
	case "qiniu":
		data = global.Config.QiNIu
	case "ai":
		data = global.Config.Ai
	default:
		res.FailWithMsg("不存在这个错误", c)
		return
	}

	res.SuccessWithData(data, c)
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
	fmt.Println("req:", req)
	res.SuccessWithMsg("更新成功", c)
	return
}
