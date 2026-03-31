package article_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum/relationship_enum"
	"blogx_server/service/focus_service"

	"github.com/gin-gonic/gin"
)

type AuthRecommentResponse struct {
	UserID       uint   `json:"userID"`
	UserNickname string `json:"userNickname"`
	UserAvatar   string `json:"userAvatar"`
	UserAbstract string `json:"userAbstract"`
}

func (ArticleApi) ArticleRecommentView(c *gin.Context) {
	cr := middlerware.GetBind[common.PageInfo](c)
	var count int
	var userIDList []uint
	global.Db.Model(&models.ArticleModel{}).Group("user_id").Select("count(*)").Scan(&count)
	global.Db.Model(&models.ArticleModel{}).Group("user_id").
		Offset(cr.GetOffset()).
		Limit(cr.GetLimit()).
		Select("user_id").
		Scan(&userIDList)

	claims, err := jwts.ParseTokenByGin(c)
	if err == nil && claims != nil {
		m := focus_service.CalcUserPatchRelationship(claims.Claims.UserID, userIDList)
		userIDList = []uint{}
		for u, relation := range m {
			if relation == relationship_enum.RelationStranger || relation == relationship_enum.RelationFans {
				userIDList = append(userIDList, u)
			}
		}
	}
	var userList []models.UserModel
	global.Db.Find(&userList, "id in ?", userIDList)
	var list = make([]AuthRecommentResponse, 0)
	for _, model := range userList {
		list = append(list, AuthRecommentResponse{
			UserID:       model.ID,
			UserNickname: model.Nickname,
			UserAvatar:   model.Avatar,
			UserAbstract: model.Abstract,
		})
	}
	res.SuccessWithList(list, count, c)
}
