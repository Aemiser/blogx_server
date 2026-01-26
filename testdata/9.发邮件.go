package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"fmt"
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func testHTMLEmail() {
	em := global.Config.Email

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", em.SendNickname, em.SendEmail)
	e.To = []string{"1530192616@qq.com"}

	e.Subject = "【测试】发送HTML邮件"
	// 主要就是在这里添加HTML字节数组内容
	e.HTML = []byte(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>邮箱验证码 - 账号注册</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
        }
        
        body {
            background-color: #fef9f5;
            color: #333;
            line-height: 1.6;
            padding: 20px;
        }
        
        .email-container {
            max-width: 520px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 4px;
            overflow: hidden;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
            border: 1px solid #f0e6dd;
        }
        
        .email-header {
            background-color: #f8f1ea;
            color: #8c5c3c;
            padding: 28px 20px;
            text-align: center;
            border-bottom: 1px solid #f0e6dd;
        }
        
        .logo {
            font-size: 22px;
            font-weight: 500;
            margin-bottom: 8px;
            color: #8c5c3c;
        }
        
        .email-title {
            font-size: 17px;
            font-weight: 400;
            color: #a67c5b;
        }
        
        .email-body {
            padding: 32px;
        }
        
        .greeting {
            font-size: 16px;
            margin-bottom: 24px;
            color: #5c4a3c;
        }
        
        .info-box {
            background-color: #fdf8f4;
            border: 1px solid #f0e6dd;
            padding: 18px;
            margin-bottom: 24px;
            font-size: 14px;
            color: #7d6452;
            border-radius: 3px;
        }
        
        .info-box p {
            margin-bottom: 6px;
        }
        
        .verification-code {
            text-align: center;
            margin: 32px 0;
            padding: 20px 0;
        }
        
        .code-label {
            font-size: 14px;
            color: #a67c5b;
            margin-bottom: 16px;
            letter-spacing: 0.5px;
        }
        
        .code {
            font-size: 38px;
            font-weight: 400;
            letter-spacing: 6px;
            color: #8c5c3c;
            padding: 18px 30px;
            display: inline-block;
            margin: 10px 0;
            background-color: #fdf8f4;
            border-radius: 4px;
            font-family: 'SF Mono', Monaco, 'Courier New', monospace;
            border: 2px solid #f0e6dd;
        }
        
        .expiry-note {
            color: #b89a83;
            margin-top: 12px;
            font-size: 13px;
        }
        
        .action-button {
            display: block;
            width: 100%;
            background-color: #d4a574;
            color: white;
            text-decoration: none;
            padding: 15px;
            border-radius: 3px;
            font-weight: 500;
            font-size: 15px;
            margin-top: 24px;
            text-align: center;
            transition: background-color 0.2s;
            border: none;
            cursor: pointer;
            letter-spacing: 0.5px;
        }
        
        .action-button:hover {
            background-color: #c49566;
        }
        
        .instructions {
            padding: 22px 0;
            margin: 24px 0;
            font-size: 13px;
            color: #7d6452;
            border-top: 1px solid #f0e6dd;
            border-bottom: 1px solid #f0e6dd;
        }
        
        .instructions h3 {
            color: #8c5c3c;
            margin-bottom: 14px;
            font-size: 15px;
            font-weight: 500;
        }
        
        .instructions ul {
            padding-left: 20px;
        }
        
        .instructions li {
            margin-bottom: 8px;
            line-height: 1.5;
        }
        
        .highlight {
            color: #8c5c3c;
            font-weight: 500;
        }
        
        .email-footer {
            background-color: #f8f1ea;
            padding: 24px;
            text-align: center;
            color: #a67c5b;
            font-size: 12px;
            border-top: 1px solid #f0e6dd;
        }
        
        .footer-links {
            margin-top: 16px;
        }
        
        .footer-links a {
            color: #b89a83;
            text-decoration: none;
            margin: 0 10px;
            font-size: 12px;
        }
        
        .footer-links a:hover {
            color: #8c5c3c;
        }
        
        .security-note {
            margin-top: 24px;
            padding: 16px;
            background-color: #fdf8f4;
            font-size: 13px;
            color: #7d6452;
            text-align: center;
            border-radius: 3px;
            border-left: 3px solid #d4a574;
        }
        
        .note-icon {
            color: #d4a574;
            margin-right: 6px;
        }
        
        @media (max-width: 600px) {
            .email-body {
                padding: 24px;
            }
            
            .code {
                font-size: 32px;
                letter-spacing: 5px;
                padding: 16px 24px;
            }
            
            .email-header {
                padding: 24px 16px;
            }
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="email-header">
            <div class="logo">验证中心</div>
            <h1 class="email-title">邮箱验证码</h1>
        </div>
        
        <div class="email-body">
            <p class="greeting">尊敬的用户，您好！</p>
            
            <div class="info-box">
                <p>您正在进行邮箱注册/验证操作</p>
                <p>请使用以下验证码完成验证：</p>
            </div>
            
            <div class="verification-code">
                <div class="code-label">验证码</div>
                <div class="code">3954</div>
                <div class="expiry-note">此验证码10分钟内有效</div>
            </div>
            
            <button class="action-button">验证邮箱</button>
            
            <div class="instructions">
                <h3>使用说明</h3>
                <ul>
                    <li>此验证码用于验证您的邮箱所有权，请勿泄露给他人</li>
                    <li>验证码有效期为 <span class="highlight">10分钟</span>，请尽快使用</li>
                    <li>如果您没有进行此操作，请忽略此邮件</li>
                    <li>请勿回复此邮件，如有问题请联系客服</li>
                </ul>
            </div>
            
            <div class="security-note">
                <span class="note-icon">•</span>为保障您的账号安全，请勿将验证码提供给任何人。
            </div>
        </div>
        
        <div class="email-footer">
            <p>此邮件由系统自动发送，请勿直接回复</p>
            <p>© 2023 验证中心 版权所有</p>
            <div class="footer-links">
                <a href="#">隐私政策</a>
                <a href="#">使用条款</a>
                <a href="#">帮助中心</a>
            </div>
        </div>
    </div>
</body>
</html>`)

	addr := fmt.Sprintf("%s:%d", em.Domain, em.Port)
	err := e.Send(addr, smtp.PlainAuth("", em.SendEmail, em.AuthCode, em.Domain))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	testHTMLEmail()
}
