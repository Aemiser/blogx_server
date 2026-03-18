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
	"blogx_server/service/redis_service/redis_comment"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (CommentApi) CommentRemoveView(c *gin.Context) {
	cr := middlerware.GetBind[models.IDRequest](c)

	var model models.CommentModel
	err := global.Db.Preload("ArticleModel").Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("评论不存在", c)
		return
	}

	claims := jwts.GetClaimsByGin(c)
	if claims.Claims.Role != enum.AdminRole {
		// 不是自己发的评论  不是自己发的文章的评论
		if model.UserID != claims.Claims.UserID || model.ArticleModel.UserID != claims.Claims.UserID {
			res.FailWithMsg("无权限", c)
			return
		}
	}
	message_service.InsertSystemMessage(model.UserID, "管理员删除了你的评论", fmt.Sprintf("【%s】 评论不符合社区规范", model.Content), "", "")

	//删评论
	// 要找到所有的子评论和所有的父评论

	//获取该评论下的所有子评论
	subList := comment_service.GetCommentOneDimensionalization(model.ID)

	// 获取该评论的父评论
	if model.ParentID != nil {
		parentList := comment_service.GetParents(*model.ParentID)
		// 对每个父评论的回复 - 去 len(subList)
		for _, subModel := range parentList {
			redis_comment.SetCacheApply(subModel.ID, -len(subList))
		}
	}

	// 删评论
	global.Db.Delete(&subList)
	res.SuccessWithMsgf(c, "删除评论成，共删除 %d 条", len(subList))
}
