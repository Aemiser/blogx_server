package global_notification_api

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

type GlobalNotificationApi struct {
}

type CreateRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Icon    string `json:"icon"`
	Href    string `json:"href"`
}

func (GlobalNotificationApi) CreateView(c *gin.Context) {
	cr := middlerware.GetBind[CreateRequest](c)

	var model models.GlobalNotificationModel
	err := global.Db.Take(&model, "title = ?", cr.Title).Error
	if err == nil {
		res.FailWithMsg("全局消息名称重复", c)
		return
	}

	err = global.Db.Create(&models.GlobalNotificationModel{
		Title:   cr.Title,
		Icon:    cr.Icon,
		Content: cr.Content,
		Href:    cr.Href,
	}).Error

	if err != nil {
		res.FailWithMsg("全局消息创建失败", c)
		return
	}
	res.SuccessWithMsg("全局消息创建成功", c)

}

type ListRequest struct {
	common.PageInfo
	Type int8 `form:"type" binding:"required,oneof=1 2"` // 1是查自己的 2 是后台
}

type ListResponse struct {
	models.GlobalNotificationModel
	IsRead bool `json:"isRead"`
}

func (GlobalNotificationApi) ListView(c *gin.Context) {
	cr := middlerware.GetBind[ListRequest](c)

	claims := jwts.GetClaimsByGin(c)
	readMsgMap := map[uint]bool{}

	// 高级查询
	query := global.Db.Where("")
	switch cr.Type {
	case 1: // 用户可见
		// 没有被用户删除的
		var ugnList []models.UserGlobalnotificationModel
		global.Db.Find(&ugnList, "user_id = ? ", claims.Claims.UserID)
		var msgIDList []uint
		for _, model := range ugnList {
			// 把删除的ID放入列表中
			if model.IsDelete {
				msgIDList = append(msgIDList, model.ID)
				continue
			}
			if model.IsRead {
				readMsgMap[model.NotificationID] = true
			}
		}
		// 如果没有就不加这句
		if len(msgIDList) > 0 {
			query = query.Where("id not in ?", msgIDList)
		}

	case 2:
		if claims.Claims.Role != enum.AdminRole {
			res.FailWithMsg("无权限", c)
			return
		}
	}
	_list, count, _ := common.ListQuery(models.GlobalNotificationModel{}, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"title", "content"},
		Where:    query,
	})

	var list []ListResponse
	for _, model := range _list {
		list = append(list, ListResponse{
			GlobalNotificationModel: model,
			IsRead:                  readMsgMap[model.ID],
		})
	}

	res.SuccessWithList(list, count, c)
}

func (GlobalNotificationApi) RemoveView(c *gin.Context) {
	cr := middlerware.GetBind[models.IDListRequest](c)

	// 查库是否存在
	var list []models.GlobalNotificationModel
	global.Db.Find(&list, "id in  ?", cr.IDList)

	if len(list) > 0 {
		global.Db.Delete(&list)
	}

	res.SuccessWithMsgf(c, "删除%d条全局消息，成功%d个", len(cr.IDList), len(list))
	return
}

type UserMsgActionRequest struct {
	ID   uint `json:"id" binding:"required"`
	Type int8 `json:"type" binding:"required,oneof=1 2"` // 1 读取 2 删除
}

// 用户读取或者删除全局消息
func (GlobalNotificationApi) UserMsgActionView(c *gin.Context) {
	cr := middlerware.GetBind[UserMsgActionRequest](c)

	var msg models.GlobalNotificationModel
	err := global.Db.Take(&msg, cr.ID).Error
	if err != nil {
		res.FailWithMsg("全局消息不存在", c)
		return
	}

	userID := jwts.GetUserIDByGin(c)

	model := models.UserGlobalnotificationModel{
		UserID:         userID,
		NotificationID: cr.ID,
	}

	m := "消息读取成功"
	if cr.Type == 1 {
		model.IsRead = true
	} else {
		model.IsDelete = true
		m = "消息删除成功"
	}
	// 看一看之前有没有操作过
	var ugn models.UserGlobalnotificationModel
	err = global.Db.Take(&ugn, "user_id = ? and notification_iD = ?", userID, cr.ID).Error
	// 之前这个用户对这个消息没有操作过
	// 之前对这个消息有读取操作
	// 之前对这个消息有删除操作
	// 先读取在删除
	if err != nil {
		global.Db.Create(&model)
		res.SuccessWithMsg("消息读取成功", c)
		return
	}

	if ugn.IsDelete {
		res.FailWithMsg("消息已删除", c)
		return
	}

	if ugn.IsRead {
		// 如果现在是删除操作，那就更新
		if model.IsDelete {
			global.Db.Model(&ugn).Update("is_delete", true)
			res.SuccessWithMsg("消息删除成功", c)
			return
		}

	}

	res.SuccessWithMsg(m, c)
}
