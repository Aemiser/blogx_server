package comment_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/models/enum/relationship_enum"
	"blogx_server/service/focus_service"
	"blogx_server/service/redis_service/redis_comment"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentListRequest struct {
	common.PageInfo
	ArticleID uint `form:"articleID"`
	UserID    uint `form:"userID"`
	Type      int8 `form:"type" binding:"required,oneof=1 2 3"` // 1查我发文章的评论 2 查我发布的评论 3 管理员
}

type CommentListResponse struct {
	ID           uint                       `json:"id"`
	CreatedAt    time.Time                  `json:"createdAt"`
	Content      string                     `json:"content"`
	UserID       uint                       `json:"userID"`
	UserNickname string                     `json:"nickName"`
	UserAvatar   string                     `json:"userAvatar"`
	ArticleID    uint                       `json:"articleID"`
	ArticleTitle string                     `json:"articleTitle"`
	ArticleCover string                     `json:"articleCover"`
	DiggCount    int                        `json:"diggCount"`
	Relation     relationship_enum.Relation `json:"relation"`
	IsMe         bool                       `json:"isMe"`
}

func (CommentApi) CommentListView(c *gin.Context) {
	cr := middlerware.GetBind[CommentListRequest](c)
	query := global.Db.Where("")
	claims := jwts.GetClaimsByGin(c)
	switch cr.Type {
	case 1: // 查我发的文章的评论
		var articleIDList []uint
		global.Db.Model(models.ArticleModel{}).
			Where("user_id = ? and status = ?", claims.Claims.UserID, enum.ArticlePublished).
			Select("id").Scan(&articleIDList)
		query.Where("article_id in ?", articleIDList)
	case 2: // 查我发布的评论
		cr.UserID = claims.Claims.UserID

	case 3: // 管理员
	}
	_list, count, _ := common.ListQuery(models.CommentModel{
		ArticleID: cr.ArticleID,
		UserID:    cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"content"},
		Preloads: []string{"UserModel", "ArticleModel"},
		Where:    query,
	})

	var RelationMap = map[uint]relationship_enum.Relation{}
	if cr.Type == 1 {
		var userIDList []uint
		for _, model := range _list {
			userIDList = append(userIDList, model.UserID)
		}
		RelationMap = focus_service.CalcUserPatchRelationship(claims.Claims.UserID, userIDList)
	}

	var list = make([]CommentListResponse, 0)
	for _, model := range _list {
		list = append(list, CommentListResponse{
			ID:           model.ID,
			CreatedAt:    model.CreatedAt,
			Content:      model.Content,
			UserID:       model.UserID,
			UserNickname: model.UserModel.Nickname,
			UserAvatar:   model.UserModel.Avatar,
			ArticleID:    model.ArticleID,
			ArticleTitle: model.ArticleModel.Title,
			ArticleCover: model.ArticleModel.Cover,
			DiggCount:    model.DiggCount + redis_comment.GetCacheDigg(model.ID),
			Relation:     RelationMap[model.UserID],
			IsMe:         model.UserID == claims.Claims.UserID,
		})
	}

	res.SuccessWithList(list, count, c)
}
