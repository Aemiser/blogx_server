package comment_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/comment_service"

	"github.com/gin-gonic/gin"
)

func (CommentApi) CommentTreeView(c *gin.Context) {
	cr := middlerware.GetBind[models.IDRequest](c)

	var article models.ArticleModel
	err := global.Db.Take(&article, "id = ? and status =?", cr.ID, enum.ArticlePublished).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	//查找跟评论
	var commentList []models.CommentModel
	err = global.Db.Find(&commentList, "article_id = ? and parent_id is null", cr.ID).Error
	var list = make([]comment_service.CommentResponse, 0)
	for _, model := range commentList {
		resp := comment_service.GetCommentTreeV4(model.ID)
		list = append(list, *resp)
	}

	res.SuccessWithList(list, len(list), c)

}
