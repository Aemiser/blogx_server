package comment_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/service/redis_service/redis_comment"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (CommentApi) CommentDiggView(c *gin.Context) {
	cr := middlerware.GetBind[models.IDRequest](c)

	var comment models.CommentModel
	err := global.Db.Take(&comment, cr.ID).Error
	if err != nil {
		res.FailWithMsg("评论不存在", c)
		return
	}

	claims := jwts.GetClaimsByGin(c)

	var digg models.CommentDiggModel
	// 使用 Unscoped 来查询包括已软删除的记录
	err = global.Db.Unscoped().Take(&digg, "comment_id = ? and user_id = ?", cr.ID, claims.Claims.UserID).Error
	if err != nil {
		// 记录不存在，创建新的点赞
		err = global.Db.Create(&models.CommentDiggModel{
			CommentID: cr.ID,
			UserID:    claims.Claims.UserID,
		}).Error
		if err != nil {
			res.FailWithMsg("点赞失败", c)
			return
		}
		res.SuccessWithMsg("点赞成功", c)
		// 写入缓存中
		redis_comment.SetCacheDigg(comment.ID, 1)

		return
	}

	// TODO :添加点赞记录到缓存中
	// 记录已存在
	if !digg.DeletedAt.Valid { // true ==>有值
		// 未删除状态，执行软删除（取消点赞）
		// GORM 会自动处理 DeletedAt 的设置
		err = global.Db.Where("comment_id = ? AND user_id = ?", cr.ID, claims.Claims.UserID).Delete(&models.CommentDiggModel{}).Error
		if err != nil {
			res.FailWithMsg("取消点赞失败", c)
			return
		}
		redis_comment.SetCacheDigg(comment.ID, -1)

		res.SuccessWithMsg("取消点赞成功", c)
	} else {
		// 已删除状态（DeletedAt 有值），恢复点赞
		digg.DeletedAt = gorm.DeletedAt{}
		err = global.Db.Where("comment_id = ? AND user_id = ?", cr.ID, claims.Claims.UserID).Save(&digg).Error
		if err != nil {
			res.FailWithMsg("点赞失败", c)
			return
		}
		res.SuccessWithMsg("点赞成功", c)
		redis_comment.SetCacheDigg(comment.ID, 1)
	}
	return
}
