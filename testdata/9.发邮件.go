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
	e.HTML = []byte(`
<!DOCTYPE html>
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
            font-family: 'Segoe UI', 'Microsoft YaHei', sans-serif;
        }
        
        body {
            background-color: #f5f7fa;
            color: #333;
            line-height: 1.6;
            padding: 20px;
        }
        
        .email-container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 5px 20px rgba(0, 0, 0, 0.08);
        }
        
        .email-header {
            background: linear-gradient(135deg, #6a11cb 0%, #2575fc 100%);
            color: white;
            padding: 30px 20px;
            text-align: center;
        }
        
        .logo {
            font-size: 28px;
            font-weight: 700;
            margin-bottom: 10px;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 10px;
        }
        
        .logo-icon {
            font-size: 32px;
        }
        
        .email-title {
            font-size: 22px;
            margin-top: 10px;
            font-weight: 600;
        }
        
        .email-body {
            padding: 30px;
        }
        
        .greeting {
            font-size: 18px;
            margin-bottom: 25px;
            color: #444;
        }
        
        .info-box {
            background-color: #f8f9ff;
            border-left: 4px solid #2575fc;
            padding: 18px 20px;
            margin-bottom: 25px;
            border-radius: 0 8px 8px 0;
        }
        
        .info-box p {
            margin-bottom: 8px;
        }
        
        .verification-code {
            text-align: center;
            margin: 35px 0;
            padding: 20px;
        }
        
        .code-label {
            font-size: 16px;
            color: #666;
            margin-bottom: 15px;
        }
        
        .code {
            font-size: 42px;
            font-weight: 700;
            letter-spacing: 8px;
            color: #2575fc;
            background-color: #f0f5ff;
            padding: 15px 25px;
            border-radius: 10px;
            display: inline-block;
            margin: 10px 0;
            border: 2px dashed #c2d6ff;
        }
        
        .instructions {
            background-color: #f9f9f9;
            padding: 20px;
            border-radius: 8px;
            margin: 25px 0;
            font-size: 14px;
            color: #666;
        }
        
        .instructions h3 {
            color: #333;
            margin-bottom: 10px;
            font-size: 16px;
        }
        
        .instructions ul {
            padding-left: 20px;
        }
        
        .instructions li {
            margin-bottom: 8px;
        }
        
        .action-button {
            display: inline-block;
            background: linear-gradient(to right, #2575fc, #6a11cb);
            color: white;
            text-decoration: none;
            padding: 14px 32px;
            border-radius: 30px;
            font-weight: 600;
            font-size: 16px;
            margin-top: 15px;
            text-align: center;
            transition: all 0.3s ease;
            box-shadow: 0 4px 12px rgba(37, 117, 252, 0.3);
        }
        
        .action-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 16px rgba(37, 117, 252, 0.4);
        }
        
        .email-footer {
            background-color: #f8f9fa;
            padding: 25px;
            text-align: center;
            color: #888;
            font-size: 13px;
            border-top: 1px solid #eee;
        }
        
        .footer-links {
            margin-top: 15px;
        }
        
        .footer-links a {
            color: #2575fc;
            text-decoration: none;
            margin: 0 10px;
        }
        
        .footer-links a:hover {
            text-decoration: underline;
        }
        
        .highlight {
            color: #2575fc;
            font-weight: 600;
        }
        
        @media (max-width: 600px) {
            .email-body {
                padding: 20px;
            }
            
            .code {
                font-size: 32px;
                letter-spacing: 5px;
                padding: 12px 18px;
            }
            
            .email-header {
                padding: 25px 15px;
            }
        }
    </style>
</head>
<body>
    <div class="email-container">
        <div class="email-header">
            <div class="logo">
                <span class="logo-icon">✉️</span>
                <span>验证中心</span>
            </div>
            <h1 class="email-title">邮箱验证码</h1>
        </div>
        
        <div class="email-body">
            <p class="greeting">尊敬的用户，您好！</p>
            
            <div class="info-box">
                <p>您正在进行邮箱注册/验证操作</p>
                <p>请使用以下验证码完成验证：</p>
            </div>
            
            <div class="verification-code">
                <div class="code-label">验证码为：</div>
                <div class="code">3954</div>
                <p style="color: #ff6b6b; margin-top: 10px; font-size: 14px;">(此验证码10分钟内有效)</p>
            </div>

            
            <div class="instructions">
                <h3>使用说明：</h3>
                <ul>
                    <li>此验证码用于验证您的邮箱所有权，请勿泄露给他人</li>
                    <li>验证码有效期为 <span class="highlight">10分钟</span>，请尽快使用</li>
                    <li>如果您没有进行此操作，请忽略此邮件</li>
                    <li>请勿回复此邮件，如有问题请联系客服</li>
                </ul>
            </div>
            
            <p style="margin-top: 25px; color: #666;">
                为保障您的账号安全，请勿将验证码提供给任何人，包括我们的工作人员。
            </p>
        </div>
        
        <div class="email-footer">
            <p>此邮件由系统自动发送，请勿直接回复</p>
            <p>© 2023 验证中心 版权所有</p>
            <div class="footer-links">
                <a href="#">隐私政策</a> | 
                <a href="#">使用条款</a> | 
                <a href="#">帮助中心</a>
            </div>
        </div>
    </div>
</body>
</html>
  `)

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
