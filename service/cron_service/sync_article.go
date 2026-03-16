package cron_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/redis_service/redis_article"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SyncArticle() {
	collectMap := redis_article.GetAllCacheCollect()
	diggMap := redis_article.GetAllCacheDigg()
	lookMap := redis_article.GetAllCacheLook()
	commentMap := redis_article.GetAllCacheComment()

	var list = []models.ArticleModel{}
	global.Db.Find(&list)

	for _, models := range list {
		collect := collectMap[models.ID]
		digg := diggMap[models.ID]
		look := lookMap[models.ID]
		comment := commentMap[models.ID]
		if collect == 0 || digg == 0 || look == 0 {
			continue
		}

		err := global.Db.Model(&models).Updates(map[string]any{
			"collect_count": gorm.Expr("collect_count + ?", collect),
			"digg_count":    gorm.Expr("digg_count + ?", digg),
			"look_count":    gorm.Expr("look_count + ?", look),
			"comment_count": gorm.Expr("lcomment_count + ?", comment),
		}).Error
		if err != nil {
			logrus.Error(err)
		}
		logrus.Infof(" %d 更新成功:", models.ID)
	}

	// 在同步的时候可能产生了一些数据
	// 可以这里在获取一遍redis里的数据，对于小数字可以这样做，但最好还是加锁
	// collectMap := redis_article.GetAllCacheCollect()
	// diggMap := redis_article.GetAllCacheDigg()
	// lookMap := redis_article.GetAllCacheLook()
	// 走完之后要清除掉
	redis_article.Clean()
}
