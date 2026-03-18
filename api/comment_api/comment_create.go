package comment_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/comment_service"
	"blogx_server/service/message_service"
	"blogx_server/service/redis_service/redis_article"
	"blogx_server/service/redis_service/redis_comment"

	"github.com/gin-gonic/gin"
)

type CommentCreateReaquest struct {
	Content   string `json:"content" binding:"required"`
	ArticleID uint   `json:"articleID" binding:"required"`
	ParentID  *uint  `json:"parentID"`
}

func (CommentApi) CommentCreateView(c *gin.Context) {
	cr := middlerware.GetBind[CommentCreateReaquest](c)

	var article models.ArticleModel
	err := global.Db.Take(&article, "id=? and status = ?", cr.ArticleID, enum.ArticlePublished).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	claims := jwts.GetClaimsByGin(c)

	model := models.CommentModel{
		Content:   cr.Content,
		UserID:    claims.Claims.UserID,
		ArticleID: cr.ArticleID,
		ParentID:  cr.ParentID,
	}
	// 找根评论
	if cr.ParentID != nil {
		// 找父评论
		parentList := comment_service.GetParents(*cr.ParentID)
		// 判断父评论的层级是否满足
		if len(parentList) >= global.Config.Site.Article.Commentline {
			res.FailWithMsg("评论层级达到限制", c)
			return
		}
		if len(parentList) > 0 {
			model.RootParentID = &parentList[len(parentList)-1].ID
			redis_comment.SetCacheApply(model.ArticleID, 1)
		}
	}

	err = global.Db.Create(&model).Error
	if err != nil {
		res.FailWithMsg("评论失败", c)
		return
	}

	redis_article.SetCacheComment(cr.ArticleID, 1)
	message_service.InsertCommentMessage(model)
	res.SuccessWithMsg("评论成功", c)
	return
}
