package chat_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"

	"github.com/gin-gonic/gin"
)

type ChatApi struct {
}

type ChatListRequest struct {
	common.PageInfo
	SendUserID uint `form:"sendUserID"`                        // 查找我和他的聊天记录
	RevUserID  uint `form:"recUserID" binding:"required"`      // 查找我和他的聊天记录
	Type       int8 `form:"type,oneof=1 2" binding:"required"` // 1：前台用户  2：管理员
}

type ChatListResponse struct {
	models.ChatModel
	SendUserNickname string `json:"sendUserNickname"`
	SendUserAvatar   string `json:"sendUserAvatar"`
	RevUserNickname  string `json:"revUserNickname"`
	RevUserAvatar    string `json:"revUserAvatar"`
	IsMe             bool   `json:"isMe"` // 是否是我发送的
	IsRead           bool   `json:"isRead"`
}

func (ChatApi) ChatListView(c *gin.Context) {
	cr := middlerware.GetBind[ChatListRequest](c)

	claims := jwts.GetClaimsByGin(c)
	userID := claims.Claims.UserID
	var deletedIDList []uint
	switch cr.Type {
	case 1: // 前台用户看
		cr.SendUserID = userID
		// 找我删除的记录

		global.Db.Model(models.UserChatAtionModel{}).Where("user_id = ? and is_delete = ?",
			userID, true).Select("chat_id").Scan(&deletedIDList)

	case 2: // 管理员看
		if claims.Claims.Role != enum.AdminRole {
			res.FailWithMsg("无权限", c)
			return
		}
		if cr.SendUserID == 0 {
			res.FailWithMsg("SendUserID必填", c)
			return
		}
	}
	query := global.Db.Where("(seed_user_id = ? and rev_user_id = ?) or (seed_user_id = ? and rev_user_id = ?)",
		cr.SendUserID, cr.RevUserID, cr.RevUserID, cr.SendUserID)

	if len(deletedIDList) > 0 {
		query = query.Where("id not in ?", deletedIDList)
	}
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
