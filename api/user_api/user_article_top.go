package user_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"

	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserArticleTopRequest struct {
	UserID    uint `json:"userID" `
	ArticleID uint `json:"articleID" binding:"required"`
	Type      int8 `json:"type" binding:"required,oneof=1 2"`
}

func (UserApi) UserArticleTopView(c *gin.Context) {
	cr := middlerware.GetBind[UserArticleTopRequest](c)

	var model models.ArticleModel
	err := global.Db.Take(&model, cr.ArticleID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	claims := jwts.GetClaimsByGin(c)
	userID := claims.Claims.UserID

	switch cr.Type {
	case 1:
		// 用户置顶文章
		// 用户只能置顶自己的文章
		if model.UserID != userID {
			res.FailWithMsg("无权限", c)
			return
		}

		// 只能置顶已发布的文章
		if model.Status != enum.ArticlePublished {
			res.FailWithMsg("文章未发布", c)
			return
		}

		// 判断之前有没有置顶过（包括已软删除的）
		var userTopArticleModel models.UserTopArticleModel
		err = global.Db.Unscoped().Where("user_id = ?", userID).First(&userTopArticleModel).Error

		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			// 之前没有置顶记录，创建新的
			global.Db.Create(&models.UserTopArticleModel{
				UserID:    userID,
				ArticleID: cr.ArticleID,
			})
			res.SuccessWithMsg("置顶文章成功", c)
		} else if userTopArticleModel.DeletedAt.Valid {
			// 有已软删除的记录，恢复它并更新文章 ID
			userTopArticleModel.DeletedAt = gorm.DeletedAt{}
			userTopArticleModel.ArticleID = cr.ArticleID
			global.Db.Save(&userTopArticleModel)
			res.SuccessWithMsg("置顶文章成功", c)
		} else {
			// 有未删除的记录，取消置顶
			global.Db.Delete(&userTopArticleModel)
			res.SuccessWithMsg("取消置顶成功", c)
		}

		return

	case 2:
		// 管理员置顶文章
		if claims.Claims.Role != enum.AdminRole {
			res.FailWithMsg("无权限", c)
			return
		}
		if model.Status != enum.ArticlePublished {
			res.FailWithMsg("管理员只能置顶已发布的文章", c)
			return
		}
		var userTopArticle models.UserTopArticleModel
		err = global.Db.Unscoped().Where("user_id = ? and article_id = ?", userID, cr.ArticleID).First(&userTopArticle).Error
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) { // 查不到就置顶
			global.Db.Create(&models.UserTopArticleModel{
				UserID:    userID,
				ArticleID: cr.ArticleID,
			})
			res.SuccessWithMsg("置顶文章成功", c)
			return
		}

		// 查到了
		global.Db.Delete(&userTopArticle)
		res.FailWithMsg("取消置顶成功", c)
	}
}
