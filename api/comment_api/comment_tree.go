package comment_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/models/enum/relationship_enum"
	"blogx_server/service/comment_service"
	"blogx_server/service/focus_service"
	"blogx_server/utils"
	"fmt"

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
	fmt.Println("articleID", article.ID)
	var userRelationMap = map[uint]relationship_enum.Relation{}
	var userDiggCommentMap = map[uint]bool{}
	claims, err := jwts.ParseTokenByGin(c)
	if err == nil && claims != nil {
		// 登入了
		var commentList []models.CommentModel //文章的评论id列表
		global.Db.Find(&commentList, "article_id = ?", cr.ID)

		if len(commentList) > 0 {
			var commentIDList []uint
			var userIDList []uint
			for _, model := range commentList {
				commentIDList = append(commentIDList, model.ID)
				userIDList = append(userIDList, model.UserID)
			}
			userIDList = utils.Unique(userIDList) // 对用户ID列表去重
			userRelationMap = focus_service.CalcUserPatchRelationship(claims.Claims.UserID, userIDList)
			var commentDiggList []models.CommentDiggModel
			global.Db.Order("created_at desc").Find(&commentDiggList, "user_id = ? and comment_id in ?", claims.Claims.UserID, commentIDList)
			for _, model := range commentDiggList {
				userDiggCommentMap[model.CommentID] = true
			}
		}
	}

	//查找跟评论
	var commentList []models.CommentModel
	err = global.Db.Find(&commentList, "article_id = ? and parent_id is null", cr.ID).Error
	var list = make([]comment_service.CommentResponse, 0)
	for _, model := range commentList {
		resp := comment_service.GetCommentTreeV4(model.ID, userRelationMap, userDiggCommentMap)
		list = append(list, *resp)
	}

	res.SuccessWithList(list, len(list), c)

}
