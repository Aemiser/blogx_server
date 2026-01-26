package captcha_api

import (
	"blogx_server/common/res"
	"blogx_server/utils"

	"github.com/gin-gonic/gin"
)

type CaptchaApi struct {
}
type CaptchaResponse struct {
	CaptchaID string `json:"captchaID"`
	Captcha   string `json:"captcha"`
}

func (CaptchaApi) CaptchaView(c *gin.Context) {
	id, captcha64, err := utils.GetCaptcha()
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	res.SuccessWithData(CaptchaResponse{
		CaptchaID: id,
		Captcha:   captcha64,
	}, c)
}
