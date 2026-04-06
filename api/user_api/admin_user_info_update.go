package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/log_service"
	"blogx_server/utils/maps"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AdminUserInfoUpdataRequest struct {
	UserID   uint           `json:"userID"  binding:"required"`
	Username *string        `json:"username" s-u:"username"`
	Nickname *string        `json:"nickname" s-u:"nickname"`
	Avatar   *string        `json:"avatar" s-u:"avatar"`
	Abstract *string        `json:"abstract" s-u:"abstract"`
	Role     *enum.RoleType `json:"role" s-u:"role"`
}

func (UserApi) AdminUserInfoUpdate(c *gin.Context) {
	var req AdminUserInfoUpdataRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	log := log_service.GetLog(c)
	log.SetTitle("<span style='color: #722ed1'>👑 管理员修改用户信息</span>")
	log.SetItem("目标用户ID", fmt.Sprintf("<span style='color: #ff4d4f; font-weight: bold'>%d</span>", req.UserID))

	reqJson, _ := json.Marshal(req)
	log.SetItem("📝 修改内容", fmt.Sprintf("<pre style='background: #f5f5f5; padding: 8px; border-radius: 4px; max-height: 150px; overflow-y: auto'>%s</pre>", string(reqJson)))

	// 查找用户
	var user models.UserModel
	err = global.Db.Take(&user, req.UserID).Error
	if err != nil {
		log.SetItemError("用户不存在", err)
		res.FailWithMsg("用户不存在", c)
		return
	}

	roleIcon := "👤"
	if user.Role == enum.AdminRole {
		roleIcon = "👑"
	}
	log.SetItem("📋 修改前用户信息", fmt.Sprintf(`<div style='background: #fff7e6; padding: 12px; border-radius: 4px; border-left: 3px solid #faad14'>
		<p>👤 昵称: <span style='color: #1890ff'>%s</span></p>
		<p>📧 用户名: <span style='color: #1890ff'>%s</span></p>
		<p>🔗 头像: <img src='%s' style='max-width: 50px; border-radius: 50%;'/></p>
		<p>💼 角色: %s <span style='color: %s'>%s</span></p>
		<p>📝 个性签名: %s</p>
	</div>`, user.Nickname, user.Username, user.Avatar, roleIcon, func() string {
		if user.Role == enum.AdminRole {
			return "#722ed1"
		}
		return "#8c8c8c"
	}(), user.Role, user.Abstract))

	// 结构体转map
	userMap := maps.StructToMap(req, "s-u")
	log.SetItem("🔄 更新字段", fmt.Sprintf("<span style='color: #52c41a'>%v</span>", func() []string {
		var fields []string
		for k := range userMap {
			fields = append(fields, k)
		}
		return fields
	}()))

	// 修改用户字段
	err = global.Db.Model(&user).Updates(userMap).Error
	if err != nil {
		log.SetItemError("修改失败", err)
		res.FailWithMsg("用户信息修改失败", c)
		return
	}

	// 重新查询获取修改后的信息
	global.Db.Take(&user, req.UserID)
	log.SetItem("✅ 修改后", fmt.Sprintf(`<div style='background: #e6f7ff; padding: 12px; border-radius: 4px; border-left: 3px solid #1890ff'>
		<p>👤 昵称: <span style='color: #52c41a'>%s</span></p>
		<p>💼 角色: <span style='color: #722ed1'>%s</span></p>
	</div>`, user.Nickname, user.Role))

	res.SuccessWithMsg("用户信息修改成功", c)
}
