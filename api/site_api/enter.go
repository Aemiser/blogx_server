package site_api

import (
	"blogx_server/common/res"
	"blogx_server/conf"
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/middlerware"
	"errors"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
		err = c.ShouldBindJSON(&data)
		result = data
	case "email":
		var data conf.Email
		err = c.ShouldBindJSON(&data)
		result = data
	case "qq":
		var data conf.QQ
		err = c.ShouldBindJSON(&data)
		fmt.Println(err)
		result = data
	case "qiniu":
		var data conf.QiNiu
		err = c.ShouldBindJSON(&data)
		result = data
	case "ai":
		var data conf.Ai
		err = c.ShouldBindJSON(&data)
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
		err = UpdateSite(s)
		if err != nil {
			res.FailWithError(err, c)
			return
		}
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

func UpdateSite(site conf.Site) error {
	if site.Project.Icon == "" && site.Project.WebPath == "" && site.Project.Title == "" &&
		site.Seo.Description == "" && site.Seo.Keywords == "" {
		return nil
	}

	if site.Project.WebPath == "" {
		return errors.New("请配置前端地址")
	}

	file, err := os.Open(site.Project.WebPath)
	if err != nil {
		return errors.New(fmt.Sprintf("%s 地址不存在", site.Project.WebPath))
	}
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		logrus.Errorf("goquery解析错误 %s \n", err)
		return errors.New("文件解析失败")
	}

	if site.Project.Title != "" {
		doc.Find("title").SetText(site.Project.Title)
	}

	if site.Project.Icon != "" {
		sele := doc.Find("link[rel='icon']")
		if sele.Length() > 0 {
			// 有就修改
			doc.Find("link[rel='icon']").SetAttr("href", site.Project.Icon)
		} else {
			// 没有就创建
			doc.Find("head").AppendHtml(fmt.Sprintf("<link rel='icon' href='%s'  />", site.Project.Icon))
		}
	}

	if site.Seo.Keywords != "" {
		sele := doc.Find("meta[name='keywords']")
		if sele.Length() > 0 {
			// 有就修改
			doc.Find("meta[name='keywords']").SetAttr("content", site.Seo.Keywords)
		} else {
			// 没有就创建
			doc.Find("head").AppendHtml(fmt.Sprintf(" <meta name=\"keywords\" content=\"%s\">'  />", site.Seo.Keywords))
		}
	}

	if site.Seo.Description != "" {
		sele := doc.Find("meta[name='description']")
		if sele.Length() > 0 {
			// 有就修改
			doc.Find("meta[name='description']").SetAttr("content", site.Seo.Description)
		} else {
			// 没有就创建
			doc.Find("head").AppendHtml(fmt.Sprintf(" <meta name=\"description\" content=\"%s\">'  />", site.Seo.Description))
		}
	}

	html, err := doc.Html()
	if err != nil {
		logrus.Errorf("生成html失败: %s \n", err)
		return errors.New("生成html失败")
	}

	// 修改文件
	err = os.WriteFile(site.Project.WebPath, []byte(html), 0666)
	if err != nil {
		logrus.Errorf("修改文件失败: %s \n", err)
		return errors.New("修改文件失败")
	}
	return nil
}
func (SiteApi) SiteInfoQQView(c *gin.Context) {
	res.SuccessWithData(global.Config.QQ.Url(), c)
	return
}
