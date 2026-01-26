package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/email_service"
	"blogx_server/utils"
	"blogx_server/utils/email_store"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	SendEmailTypeRegister int8 = 1
	SendEmailTypeReset    int8 = 2
)

type SendEmailRequest struct {
	Type  int8   `json:"type" binding:"oneof=1 2" `
	Email string `json:"email" binding:"required" `
}

type SendEmailResponse struct {
	EmailID string `json:"emailID"`
}

// SendEmailView 发送到邮件的验证码
func (UserApi) SendEmailView(c *gin.Context) {
	var req SendEmailRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	code := utils.GetRandomInDigital(4)
	emailId := utils.GetUUID()
	switch req.Type {
	case SendEmailTypeRegister:
		// 检查邮箱是否存在
		var model models.UserModel
		err = global.Db.Take(&model, "email = ?", req.Email).Error
		if err == nil {
			res.FailWithMsg("邮箱已存在", c)
			return
		}
		err = email_service.SendRegisteredCode(req.Email, code)
	case SendEmailTypeReset:
		err = email_service.SendRegisteredCode(req.Email, code)
	}
	if err != nil {
		logrus.Errorf("邮件发送失败：%s", err.Error())
		res.FailWithError(err, c)
		return
	}
	// 发送成功，保存captchaID
	global.EmailVerifyStore.Store(emailId, email_store.EmailStoreInfo{
		Email: req.Email,
		Code:  code,
	})
	res.SuccessWithData(SendEmailResponse{
		EmailID: emailId,
	}, c)
}
