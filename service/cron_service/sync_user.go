package cron_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/redis_service/redis_article"
	"blogx_server/service/redis_service/redis_user"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SyncUser() {
	lookMap := redis_article.GetAllCacheLook()

	var list = []models.UserConfigModel{}
	global.Db.Find(&list)

	for _, models := range list {
		look := lookMap[models.UserID]
		if look == 0 {
			continue
		}

		err := global.Db.Model(&models).Updates(map[string]any{
			"look_count": gorm.Expr("look_count + ?", look),
		}).Error
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof(" %d 更新成功:", models.UserID)
	}

	// 走完之后要清除掉
	redis_user.Clean()
}
