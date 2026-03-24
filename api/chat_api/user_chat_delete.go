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

func (ChatApi) UserChatDeleteView(c *gin.Context) {
	cr := middlerware.GetBind[models.IDListRequest](c)
	userID := jwts.GetUserIDByGin(c)

	// 查找对应的聊天信息
	var chatList []models.ChatModel
	global.Db.Find(&chatList, "id in ?", cr.IDList)

	// 之前是否操作过
	chatMap := common.ScanMapV2(models.UserChatAtionModel{}, common.ScanMapOptions{
		Where: global.Db.Where("user_id = ? and chat_id in ?", userID, cr.IDList),
		Key:   "ChatID",
	})

	var addChatAcList []models.UserChatAtionModel
	var updateChatAcIDList []uint
	// 判断消息是不是被操作过了
	for _, model := range chatList {
		chat, ok := chatMap[model.ID]
		if !ok { // 如果未存在，即创建
			addChatAcList = append(addChatAcList, models.UserChatAtionModel{
				UserID:   userID,
				ChatID:   model.ID,
				IsDelete: true,
			})
			continue
		}
		if chat.IsDelete {
			continue
		}
		updateChatAcIDList = append(updateChatAcIDList, chat.ID)
	}

	if len(addChatAcList) > 0 {
		err := global.Db.Debug().Create(&addChatAcList).Error
		if err != nil {
			res.FailWithMsg("删除消息失败", c)
			return
		}
	}

	if len(updateChatAcIDList) > 0 {
		err := global.Db.Debug().Model(&models.UserChatAtionModel{}).Where("id in ?", updateChatAcIDList).Update("is_delete", true).Error
		if err != nil {
			res.FailWithMsg("删除消息失败", c)
			return
		}
	}
	res.SuccessWithMsg("删除消息成功", c)
	return

}
