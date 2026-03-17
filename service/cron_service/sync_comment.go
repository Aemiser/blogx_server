package cron_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/redis_service/redis_comment"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SyncComment() {
	commentMap := redis_comment.GetAllCacheApply()

	var list = []models.CommentModel{}
	global.Db.Find(&list)

	for _, models := range list {
		apply := commentMap[models.ID]
		if apply == 0 {
			continue
		}

		err := global.Db.Model(&models).Updates(map[string]any{
			"apply_count": gorm.Expr("apply_count + ?", apply),
		}).Error
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof("评论 %d 更新成功:", models.ID)
	}

	// 在同步的时候可能产生了一些数据
	// 可以这里在获取一遍redis里的数据，对于小数字可以这样做，但最好还是加锁
	// collectMap := redis_article.GetAllCacheCollect()
	// diggMap := redis_article.GetAllCacheDigg()
	// lookMap := redis_article.GetAllCacheLook()
	// 走完之后要清除掉
	redis_comment.Clean()
}
