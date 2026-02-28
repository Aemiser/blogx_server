package core

import (
	"blogx_server/global"
	"blogx_server/service/river_service"

	"github.com/sirupsen/logrus"
)

func InitMysqlEs() {
	if global.Config.ES.Addr == "" {
		logrus.Infof("未配置es")
		return
	}

	if !global.Config.River.Enable {
		logrus.Infof("未启用es和mysql的同步")
		return
	}

	// 如果启用了调试模式，则跳过 river service
	if global.Config.River.Debug {
		logrus.Infof("River service 调试模式已启用，跳过数据同步")
		return
	}

	r, err := river_service.NewRiver()
	if err != nil {
		logrus.Fatalf("river_service.NewRiver err:%v", err)
	}
	go r.Run()
}
