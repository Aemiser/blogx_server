package middlerware

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

type captchaMiddlewareRequest struct {
	CaptchaID   string `json:"captchaID"`
	CaptchaCode string `json:"captchaCode"`
}

func CaptchaMiddleware(c *gin.Context) {
	if !global.Config.Site.Login.Captcha {
		return
	}

	byteData, err := c.GetRawData()
	if err != nil {
		res.FailWithMsgf(c, "获取请求体失败")
		c.Abort()
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(byteData))
	var cr captchaMiddlewareRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithMsg("获取图形验证失败", c)
		c.Abort()
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(byteData))

	// 验证图形验证码
	if !global.Stores.Verify(cr.CaptchaID, cr.CaptchaCode, true) {
		res.FailWithMsg("图形验证码错误", c)
		c.Abort()
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(byteData))
}
