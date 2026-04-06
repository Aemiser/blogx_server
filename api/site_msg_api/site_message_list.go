package site_msg_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum/message_type_enum"
	"blogx_server/models/enum/relationship_enum"
	"blogx_server/service/focus_service"

	"github.com/gin-gonic/gin"
)

type SiteMsgListRequest struct {
	common.PageInfo
	T int8 `form:"t" binding:"required,oneof=1 2 3"` // 1评论和回复 2赞和收藏 3 系统
}
type SiteMsgListResponse struct {
	models.MessageModel
	Relation relationship_enum.Relation
}

func (SiteMsgApi) SiteMsgListView(c *gin.Context) {
	cr := middlerware.GetBind[SiteMsgListRequest](c)

	claims := jwts.GetClaimsByGin(c)
	var typeList []message_type_enum.Type
	switch cr.T {
	case 1:
		typeList = append(typeList, message_type_enum.CommentType, message_type_enum.ApplyType)
	case 2:
		typeList = append(typeList, message_type_enum.DiggArticleType, message_type_enum.CollectArticleType, message_type_enum.DiggCommentType)
	case 3:
		typeList = append(typeList, message_type_enum.SystemType)
		// 全局消息
	}
	_list, count, _ := common.ListQuery(models.MessageModel{
		RecvUserID: claims.Claims.UserID,
	}, common.Options{
		PageInfo:     cr.PageInfo,
		Where:        global.Db.Where("type in ?", typeList),
		DefaultOrder: "created_at desc",
	})

	var userIDList []uint
	for _, model := range _list {
		if model.ActionUserID != 0 {
			userIDList = append(userIDList, model.ActionUserID)
		}
	}
	var m map[uint]relationship_enum.Relation
	if len(userIDList) > 0 {
		m = focus_service.CalcUserPatchRelationship(claims.Claims.UserID, userIDList)
	}

	var list = make([]SiteMsgListResponse, 0)
	for _, model := range _list {
		list = append(list, SiteMsgListResponse{
			MessageModel: model,
			Relation:     m[model.ActionUserID],
		})
	}
	res.SuccessWithList(list, count, c)
}
