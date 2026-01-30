package flags_user

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils"
	"blogx_server/utils/pwd"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

type FlagUser struct {
}

func (FlagUser) Create() {
	var role enum.RoleType
	fmt.Println("选择角色  1.管理员  2.普通角色  3.游客")
	_, err := fmt.Scan(&role)
	if err != nil {
		logrus.Errorf("输入错误 %s", err)
		return
	}

	if role < 1 || role > 3 {
		logrus.Errorf("输入角色错误")
		return
	}

	var username string
	fmt.Println("请输入用户名：")
	fmt.Scan(&username)

	var model models.UserModel
	err = global.Db.Take(&model, "username = ?", username).Error
	if err == nil {
		logrus.Errorf("用户已存在")
		return
	}

	fmt.Println("请输入密码：")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		logrus.Errorf("输入密码错误 %s", err)
		return
	}

	fmt.Println("请再次输入密码：")
	repassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		logrus.Errorf("再次输入密码错误 %s", err)
		return
	}

	if string(password) != string(repassword) {
		logrus.Errorf("两次输入的密码不一致")
		return
	}

	hashpwd, _ := pwd.GenerateHashPassword(string(password))
	nickname := fmt.Sprintf("%s%s", username, utils.GetRandomWord(6))
	//创建 用户
	err = global.Db.Create(&models.UserModel{
		Username:       username,
		Role:           role,
		Password:       hashpwd,
		Nickname:       nickname,
		RegisterSource: enum.RegisterSourceTypeCmd,
	}).Error

	if err != nil {
		logrus.Errorf("创建用户失败 %s", err)
		return
	}
	logrus.Infof("创建用户成功")
}
