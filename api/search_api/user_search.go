package search_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum/relationship_enum"
	"blogx_server/service/focus_service"

	"github.com/gin-gonic/gin"
)

type UserSearchRequest struct {
	common.PageInfo
}
type UserSearchResponse struct {
	UserID   uint   `json:"userID"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Abstract string `json:"abstract"`
	Relation relationship_enum.Relation
}

func (SearchApi) UserSearchView(c *gin.Context) {
	cr := middlerware.GetBind[UserSearchRequest](c)
	_list, count, _ := common.ListQuery(models.UserModel{}, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"nickname"},
	})

	// 获取用户ID
	var userIDList []uint
	for _, item := range _list {
		userIDList = append(userIDList, item.ID)
	}

	claims, err := jwts.ParseTokenByGin(c)
	var m map[uint]relationship_enum.Relation
	if claims != nil && err == nil {
		userID := claims.Claims.UserID
		m = focus_service.CalcUserPatchRelationship(userID, userIDList)
	}

	var list = make([]UserSearchResponse, 0)
	for _, model := range _list {
		item := UserSearchResponse{
			UserID:   model.ID,
			Nickname: model.Nickname,
			Avatar:   model.Avatar,
			Abstract: model.Abstract,
		}
		if m != nil {
			item.Relation = m[model.ID]
		}

		list = append(list, item)
	}

	res.SuccessWithList(list, count, c)
}
