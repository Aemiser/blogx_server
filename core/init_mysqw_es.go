package core

import (
	"blogx_server/global"
	"blogx_server/service/river_service"

	"github.com/sirupsen/logrus"
)

func InitMysqlEs() {
	if !global.Config.River.Enable {
		logrus.Infof("未启用es和mysql的同步")
		return
	}
	r, err := river_service.NewRiver()
	if err != nil {
		logrus.Fatalf("river_service.NewRiver err:%v", err)
	}
	go r.Run()
}
