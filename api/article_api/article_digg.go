package article_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (ArticleApi) ArticleDiggView(c *gin.Context) {
	cr := middlerware.GetBind[models.IDRequest](c)

	var article models.ArticleModel
	err := global.Db.Take(&article, cr.ID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	claims := jwts.GetClaimsByGin(c)

	var digg models.ArticleDiggModel
	// 使用 Unscoped 来查询包括已软删除的记录
	err = global.Db.Unscoped().Take(&digg, "article_id = ? and user_id = ?", cr.ID, claims.Claims.UserID).Error
	if err != nil {
		// 记录不存在，创建新的点赞
		err = global.Db.Create(&models.ArticleDiggModel{
			ArticleID: cr.ID,
			UserID:    claims.Claims.UserID,
		}).Error
		if err != nil {
			res.FailWithMsg("点赞失败", c)
			return
		}
		res.SuccessWithMsg("点赞成功", c)
		return
	}

	// TODO :添加点赞记录到缓存中
	// 记录已存在
	if !digg.DeletedAt.Valid { // true ==>有值
		// 未删除状态，执行软删除（取消点赞）
		// GORM 会自动处理 DeletedAt 的设置
		err = global.Db.Where("article_id = ? AND user_id = ?", cr.ID, claims.Claims.UserID).Delete(&models.ArticleDiggModel{}).Error
		if err != nil {
			res.FailWithMsg("取消点赞失败", c)
			return
		}
		res.SuccessWithMsg("取消点赞成功", c)
	} else {
		// 已删除状态（DeletedAt 有值），恢复点赞
		digg.DeletedAt = gorm.DeletedAt{}
		err = global.Db.Where("article_id = ? AND user_id = ?", cr.ID, claims.Claims.UserID).Save(&digg).Error
		if err != nil {
			res.FailWithMsg("点赞失败", c)
			return
		}
		res.SuccessWithMsg("点赞成功", c)
	}
}
