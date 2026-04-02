package site_msg_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum/message_type_enum"

	"github.com/gin-gonic/gin"
)

type SiteMsgReadRequest struct {
	ID uint `json:"id"` // 消息id
	T  int8 `json:"t"`  // 消息类型
}

func (SiteMsgApi) SiteMsgReadView(c *gin.Context) {
	cr := middlerware.GetBind[SiteMsgReadRequest](c)

	claims := jwts.GetClaimsByGin(c)
	if cr.ID != 0 {
		// 找这个消息是否存在
		var msg models.MessageModel
		err := global.Db.Take(&msg, "recv_user_id =?  and id = ? ", claims.Claims.UserID, cr.ID).Error
		if err != nil {
			res.FailWithMsg("消息不存在", c)
			return
		}

		if msg.IsRead {
			res.SuccessWithMsg("消息已读", c)
		}

		global.Db.Model(&msg).Update("is_read", true)
		res.SuccessWithData("消息读取成功", c)
		return
	}

	var typeList []message_type_enum.Type
	switch cr.T {
	case 1:
		typeList = append(typeList, message_type_enum.CommentType, message_type_enum.ApplyType)
	case 2:
		typeList = append(typeList, message_type_enum.DiggArticleType, message_type_enum.CollectArticleType, message_type_enum.DiggCommentType)
	case 3:
		typeList = append(typeList, message_type_enum.SystemType)
	}

	var msgList []models.MessageModel
	global.Db.Find(&msgList, "recv_user_id =?   and type in ? and is_read =?", claims.Claims.UserID, typeList, false)

	if len(msgList) > 0 {
		global.Db.Model(&msgList).Update("is_read", true)
	}

	res.SuccessWithMsgf(c, "批量读取%d消息成功", len(msgList))

}
