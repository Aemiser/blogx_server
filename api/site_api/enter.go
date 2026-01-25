package site_api

import (
	"blogx_server/common/res"
	"blogx_server/conf"
	"blogx_server/core"
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
		global.Config.Site.About.Version = global.Version
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
		result := global.Config.Email
		result.AuthCode = "******"
		data = result
	case "qq":
		result := global.Config.QQ
		result.AppKey = "******"
		data = result
	case "qiniu":
		result := global.Config.QiNIu
		result.SecretKey = "******"
		data = result
	case "ai":
		result := global.Config.Ai
		result.SecretKey = "******"
		data = result
	default:
		res.FailWithMsg("不存在这个错误", c)
		return
	}

	res.SuccessWithData(data, c)
	return
}

type SiteUpdateRequest struct {
	Name string `json:"name" uri:"name" binding:"required"`
}

func (SiteApi) SiteUpdateView(c *gin.Context) {
	var ud SiteUpdateRequest
	err := c.ShouldBindUri(&ud)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	var result any
	switch ud.Name {
	case "site":
		var data conf.Site
		err = c.ShouldBind(&data)
		result = data
	case "email":
		var data conf.Email
		err = c.ShouldBind(&data)
		result = data
	case "qq":
		var data conf.QQ
		err = c.ShouldBind(&data)
		result = data
	case "qiniu":
		var data conf.QiNiu
		err = c.ShouldBind(&data)
		result = data
	case "ai":
		var data conf.Ai
		err = c.ShouldBind(&data)
		result = data
	default:
		res.FailWithMsg("不存在这个配置", c)
		return
	}
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	switch s := result.(type) {
	case conf.Site:
		// TODO :判断前端传来的配置
		global.Config.Site = s
	case conf.Email:
		if s.AuthCode == "******" {
			s.AuthCode = global.Config.Email.AuthCode
		}
		global.Config.Email = s
	case conf.QQ:
		if s.AppKey == "******" {
			s.AppKey = global.Config.QQ.AppKey
		}
		fmt.Println("s:", s)
		global.Config.QQ = s
	case conf.QiNiu:
		if s.SecretKey == "******" {
			s.SecretKey = global.Config.QiNIu.SecretKey
		}
		global.Config.QiNIu = s
	case conf.Ai:
		if s.SecretKey == "******" {
			s.SecretKey = global.Config.Ai.SecretKey
		}
		global.Config.Ai = s
	}
	// 保存修改的配置
	core.WriteConf()
	res.SuccessWithMsgf(c, "修改%s成功", ud.Name)
	return
}

func (SiteApi) SiteInfoQQView(c *gin.Context) {
	res.SuccessWithData(global.Config.QQ.Url(), c)
	return
}
