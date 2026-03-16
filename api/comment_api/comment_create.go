package comment_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"

	"github.com/gin-gonic/gin"
)

type CommentCreateReaquest struct {
	Content   string `json:"content" binding:"required"`
	ArticleID uint   `json:"articleID" binding:"required"`
	ParentID  *uint  `json:"parentID"`
}

func (CommentApi) CommentListView(c *gin.Context) {
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
	}

	err = global.Db.Create(&model).Error
	if err != nil {
		res.FailWithMsg("评论失败", c)
		return
	}

	res.SuccessWithMsg("评论成功", c)
	return
}
