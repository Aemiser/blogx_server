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

type SiteMsgRemoveRequest struct {
	ID uint `json:"id"`
	T  int8 `json:"t"`
}

func (SiteMsgApi) SiteMsgRemoveView(c *gin.Context) {
	cr := middlerware.GetBind[SiteMsgRemoveRequest](c)

	claims := jwts.GetClaimsByGin(c)
	if cr.ID != 0 {
		// 找这个消息是否存在
		var msg models.MessageModel
		err := global.Db.Take(&msg, "recv_user_id =?  and id = ? ", claims.Claims.UserID, cr.ID).Error
		if err != nil {
			res.FailWithMsg("消息不存在", c)
			return
		}

		global.Db.Delete(&msg)
		res.SuccessWithData("消息删除成功", c)
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
	global.Db.Find(&msgList, "recv_user_id =?  and type in ? ", claims.Claims.UserID, typeList)

	if len(msgList) > 0 {
		global.Db.Delete(&msgList)
	}

	res.SuccessWithMsgf(c, "批量删除%d消息成功", len(msgList))

}
