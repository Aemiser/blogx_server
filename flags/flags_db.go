package flags

import (
	"blogx_server/global"
	"blogx_server/models"

	"github.com/sirupsen/logrus"
)

func FlagDB() {
	err := global.Db.AutoMigrate(
		&models.UserModel{},
		&models.UserConfigModel{},
		&models.ArticleModel{},
		&models.ArticleDiggModel{},
		&models.CategoryModel{},
		&models.UserArticleCollectModel{},
		&models.CollectModel{},
		&models.ImageModel{},
		&models.UserTopArticleModel{},
		&models.UserArticleLookHistoryModel{}, //用户浏览历史记录表
		&models.CommentModel{},
		&models.LogModel{},
		&models.BannerModel{},
		&models.GlobalNotificationModel{}, // 全局通知表
		&models.UserLoginModel{},          // 用户日子志表

	)
	if err != nil {
		logrus.Warnf("数据库迁移失败 %s", err)
	}
	logrus.Info("数据库迁移成功")

}
