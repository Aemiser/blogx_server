package site_msg_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/utils/maps"

	"github.com/gin-gonic/gin"
)

func (SiteMsgApi) UserSiteMessageConfView(c *gin.Context) {
	claims := jwts.GetClaimsByGin(c)

	var userMessageConf models.UserMessageConfModel
	err := global.Db.Take(&userMessageConf, "user_id = ? ", claims.Claims.UserID).Error
	if err != nil {
		res.FailWithMsg("用户消息配置不存在", c)
		return
	}

	res.SuccessWithData(userMessageConf, c)
}

type UserMessageConfModelRequest struct {
	OpenCommentMessage *bool `json:"openCommentMessage" u:"open_comment_message"` //是否开启评论消息
	OpenDiggMessage    *bool `json:"openDiggMessage" u:"open_digg_message"`       //是否开启点赞
	OpenPrivateChat    *bool `json:"openPrivateChat" u:"open_private_chat"`       //是否开启私聊
}

func (SiteMsgApi) UserSiteMessageUpdateConfView(c *gin.Context) {
	cr := middlerware.GetBind[UserMessageConfModelRequest](c)
	claims := jwts.GetClaimsByGin(c)

	var userMessageConf models.UserMessageConfModel
	err := global.Db.Take(&userMessageConf, "user_id = ? ", claims.Claims.UserID).Error
	if err != nil {
		res.FailWithMsg("用户消息配置不存在", c)
		return
	}
	mp := maps.StructToMap(cr, "u")
	global.Db.Model(&userMessageConf).Updates(mp)

	res.SuccessWithMsg("用户配置信息成功", c)
}
