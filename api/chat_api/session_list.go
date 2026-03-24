package chat_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/ctype/chat_type"
	"blogx_server/models/enum/chat_msg_type"
	"blogx_server/models/enum/relationship_enum"
	"blogx_server/service/focus_service"
	"blogx_server/utils/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

type SessionListRequest struct {
	common.PageInfo
}
type SessionTable struct {
	SU        uint   `gorm:"column:sU"`
	RU        uint   `gorm:"column:rU"`
	MaxDate   string `gorm:"column:maxDate"`
	Count     int    `gorm:"column:c"`
	NewChatID uint   `gorm:"column:newChatID"`
}

type SessionListResponse struct {
	UserID       uint                       `json:"userID"`
	UserNickname string                     `json:"userNickname"`
	UserAvatar   string                     `json:"userAvatar"`
	Msg          chat_type.ChatMsg          `json:"msg"`
	MsgType      chat_msg_type.MsgType      `json:"msgType"`
	NewMsgDate   string                     `json:"newMsgDate"`
	Relation     relationship_enum.Relation `json:"relation"` // 好友关系
}

func (ChatApi) SessionListView(c *gin.Context) {
	cr := middlerware.GetBind[SessionListRequest](c)
	claims := jwts.GetClaimsByGin(c)
	userID := claims.Claims.UserID
	var deletedIDList []uint
	global.Db.Model(models.UserChatAtionModel{}).Where("user_id = ? and is_delete = ?",
		userID, true).Select("chat_id").Scan(&deletedIDList)

	query := global.Db.Where("")
	var _list []SessionTable
	var colunmn = fmt.Sprintf("  (select id  from chat_models where ((seed_user_id = sU and rev_user_id = rU)  or (seed_user_id = rU and rev_user_id = sU)) order by created_at desc  limit 1) as newChatID")
	if len(deletedIDList) > 0 {
		query = query.Where("id not in ?", deletedIDList)
		colunmn = fmt.Sprintf("  (select id  from chat_models where ((seed_user_id = sU and rev_user_id = rU)  or (seed_user_id = rU and rev_user_id = sU)) and id not in %s order by created_at desc  limit 1) as newChatID", sql.CoverSliceSql(deletedIDList))
	}
	global.Db.Debug().Model(models.ChatModel{}).
		Select(
			"least(seed_user_id, rev_user_id)    as sU",
			"greatest(seed_user_id, rev_user_id) as rU",
			"max(created_at) as maxDate",
			"count(*) as c",
			colunmn,
		).
		Where(query).
		Where(
			"(seed_user_id = ? or rev_user_id = ?)",
			userID, userID,
		).
		Group("least(seed_user_id, rev_user_id)").
		Group("greatest(seed_user_id, rev_user_id)").
		Order("maxDate desc").
		Limit(cr.GetLimit()).Offset(cr.GetOffset()).Scan(&_list)

	var count int
	global.Db.Select("count(*)").Table("(?) as x",
		global.Db.
			Model(models.ChatModel{})).
		Select("count(*)").
		Where(query).
		Where("(seed_user_id = ? or rev_user_id = ? )", userID, userID).
		Group("least(seed_user_id, rev_user_id)").
		Group("greatest(seed_user_id, rev_user_id)").
		Scan(&count)

	var userIDList []uint
	var chatIDList []uint
	for _, table := range _list {
		chatIDList = append(chatIDList, table.NewChatID)
		if table.RU == userID {
			userIDList = append(userIDList, table.SU)
		}
		if table.SU == userID {
			userIDList = append(userIDList, table.RU)
		}
	}

	userMap := common.ScanMapV2(models.UserModel{}, common.ScanMapOptions{
		global.Db.Where("id in ?", userIDList),
	})

	chatMap := common.ScanMapV2(models.ChatModel{}, common.ScanMapOptions{
		Where: global.Db.Where("id in ?", chatIDList),
	})

	relationMap := focus_service.CalcUserPatchRelationship(userID, userIDList)

	var list = make([]SessionListResponse, 0)
	for _, table := range _list {
		item := SessionListResponse{}
		if table.RU == userID {
			item.UserID = table.SU
		}
		if table.SU == userID {
			item.UserID = table.RU
		}

		item.UserNickname = userMap[item.UserID].Nickname
		item.UserAvatar = userMap[item.UserID].Avatar
		item.Msg = chatMap[table.NewChatID].Msg
		item.MsgType = chatMap[table.NewChatID].MsgType
		item.NewMsgDate = chatMap[table.NewChatID].CreatedAt.Format("2006-01-02 15:04:05")
		item.Relation = relationMap[item.UserID]
		list = append(list, item)
	}
	res.SuccessWithList(list, count, c)

}
