package article_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/redis_service/redis_article"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleCollectRequest struct {
	ArticleID uint `json:"articleID" binding:"required"`
	CollectID uint `json:"collectID"`
}

func (ArticleApi) ArticleCollectView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleCollectRequest](c)

	var article models.ArticleModel
	// 检查文章是否存在，判断文章状态
	err := global.Db.Take(&article, "status = ? and id =?", enum.ArticlePublished, cr.ArticleID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}
	claims := jwts.GetClaimsByGin(c)

	var collectmodel models.CollectModel
	//是默认收藏夹
	if cr.CollectID == 0 {
		err = global.Db.Take(&collectmodel, "user_id = ? and is_default = ?", claims.Claims.UserID, true).Error
		if err != nil {
			// 没有默认收藏夹，创建默认收藏夹
			collectmodel.Title = "默认收藏夹"
			collectmodel.UserID = claims.Claims.UserID
			collectmodel.IsDefault = true
			global.Db.Create(&collectmodel)
			cr.CollectID = collectmodel.ID
		}
		// 获取该用户的默认收藏夹id
		cr.CollectID = collectmodel.ID
	} else {
		// 判断收藏夹是否存在，并且是否是自己创建的
		err = global.Db.Take(&collectmodel, "user_id = ?", claims.Claims.UserID).Error
		if err != nil {
			res.FailWithMsg("收藏夹不存在", c)
			return
		}
	}

	// 判断该文章是否已被该用户收藏
	var articleCollect models.UserArticleCollectModel
	// 使用查询包括已软删除的数据
	err = global.Db.Unscoped().Where(models.UserArticleCollectModel{
		ArticleID: cr.ArticleID,
		UserID:    claims.Claims.UserID,
		CollectID: cr.CollectID}).Take(&articleCollect).Error

	if err != nil { // 说明查询为空，未收藏，先创建，
		// 文章收藏
		err = global.Db.Create(&models.UserArticleCollectModel{
			ArticleID: cr.ArticleID,
			UserID:    claims.Claims.UserID,
			CollectID: cr.CollectID,
		}).Error

		if err != nil {
			res.FailWithMsg("收藏失败", c)
			return
		}
		res.FailWithMsg("文章已收藏", c)
		// TODO:加如缓存功能
		redis_article.SetCacheCollect(article.ID, true)
		global.Db.Model(&collectmodel).Update("article_count", gorm.Expr("article_count + 1"))
		return
	}

	// 说明存在记录，判断器delect 字段是否未空
	if !articleCollect.DeletedAt.Valid {
		// 未空,则删除，即取消收藏
		err = global.Db.Model(models.UserArticleCollectModel{}).Delete(&articleCollect).Error
		if err != nil {
			res.FailWithMsg("取消收藏失败", c)
			return
		}
		res.FailWithMsg("取消收藏成功", c)
		redis_article.SetCacheCollect(article.ID, false)

		global.Db.Model(&collectmodel).Update("article_count", gorm.Expr("article_count - 1"))
		return
	}

	// 不为空，则恢复
	err = global.Db.Unscoped().Model(&articleCollect).Update("deleted_at", nil).Error
	if err != nil {
		res.FailWithMsg("收藏收藏失败", c)
		return
	}
	res.FailWithMsg("文章收藏成功", c)
	// TODO:加如缓存功能
	redis_article.SetCacheCollect(article.ID, true)
	global.Db.Model(&collectmodel).Update("article_count", gorm.Expr("article_count + 1"))
	return
}

func (ArticleApi) ArticleCollectRemoveView(c *gin.Context) {
	cr := middlerware.GetBind[models.IDListRequest](c)

	claims := jwts.GetClaimsByGin(c)

	var userCollectList []models.UserArticleCollectModel
	global.Db.Find(&userCollectList, "id in ? and user_id = ?", cr.IDList, claims.Claims.UserID)

	if len(userCollectList) > 0 {
		global.Db.Delete(&userCollectList)
	}

	res.SuccessWithMsgf(c, "批量删除文章共 %d 条", len(userCollectList))
}
