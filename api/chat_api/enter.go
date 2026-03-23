package chat_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

type ChatApi struct {
}

type ChatListRequest struct {
	common.PageInfo
	UserID uint `form:"userID"` // 查找我和他的聊天记录
}

type ChatListResponse struct {
	models.ChatModel
	SendUserNickname string `json:"sendUserNickname"`
	SendUserAvatar   string `json:"sendUserAvatar"`
	RevUserNickname  string `json:"revUserNickname"`
	RevUserAvatar    string `json:"revUserAvatar"`
	IsMe             bool   `json:"isMe"` // 是否是我发送的
}

func (ChatApi) ChatListView(c *gin.Context) {
	cr := middlerware.GetBind[ChatListRequest](c)

	userID := jwts.GetUserIDByGin(c)

	query := global.Db.Where("(seed_user_id = ? and rev_user_id = ?) or (seed_user_id = ? and rev_user_id = ?)",
		cr.UserID, userID, userID, cr.UserID)
	_list, count, _ := common.ListQuery(models.ChatModel{}, common.Options{
		PageInfo: cr.PageInfo,
		Preloads: []string{"SeedUserModel", "RevUserModel"},
		Where:    query,
	})

	var list []ChatListResponse
	for _, model := range _list {
		item := ChatListResponse{
			ChatModel:        model,
			SendUserNickname: model.SeedUserModel.Nickname,
			SendUserAvatar:   model.RevUserModel.Avatar,
			RevUserNickname:  model.SeedUserModel.Nickname,
			RevUserAvatar:    model.RevUserModel.Avatar,
		}

		if model.SeedUserID == userID {
			item.IsMe = true
		}
		list = append(list, item)
	}
	res.SuccessWithList(list, count, c)
}
