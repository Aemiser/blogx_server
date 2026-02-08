package email_service

import (
	"blogx_server/global"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
)

// SendRegisteredCode 发送注册验证码
func SendRegisteredCode(to, code string) error {
	em := global.Config.Email
	subject := fmt.Sprintf("【%s】账号注册", em.SendNickname)
	contest := fmt.Sprintf("你正在进行邮箱注册,这是你的验证码：%s,(十分钟内有效)", code)
	return sendEmail(to, subject, contest)

}

// SendResetCode 重置密码发送验证码
func SendResetCode(to, code string) error {
	em := global.Config.Email
	subject := fmt.Sprintf("【%s】账号密码重置", em.SendNickname)
	contest := fmt.Sprintf("你正在进行邮箱账号密码重置,这是你的验证码：%s,(十分钟内有效)", code)
	return sendEmail(to, subject, contest)

}

// SendBindCode 绑定邮箱发送验证码
func SendBindCode(to, code string) error {
	em := global.Config.Email
	subject := fmt.Sprintf("【%s】绑定邮箱操作", em.SendNickname)
	contest := fmt.Sprintf("你正在进行绑定邮箱操作,这是你的验证码：%s,(十分钟内有效)", code)
	return sendEmail(to, subject, contest)
}

func sendEmail(to, subject, content string) (err error) {
	em := global.Config.Email

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", em.SendNickname, em.SendEmail)
	e.To = []string{to}

	e.Subject = subject
	// 主要就是在这里添加HTML字节数组内容
	e.Text = []byte(content)

	addr := fmt.Sprintf("%s:%d", em.Domain, em.Port)
	err = e.Send(addr, smtp.PlainAuth("", em.SendEmail, em.AuthCode, em.Domain))
	if err != nil {
		logrus.Warnf("发送邮件失败:%s", err)
		return
	}
	return
}
