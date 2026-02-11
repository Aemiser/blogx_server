package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/maps"

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
	// 查找用户
	var user models.UserModel
	err = global.Db.Take(&user, req.UserID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}
	// 结构体转map
	userMap := maps.StructToMap(req, "s-u")
	// 修改用户字段
	err = global.Db.Model(&user).Updates(userMap).Error
	if err != nil {
		res.FailWithMsg("用户信息修改失败", c)
		return
	}
	res.SuccessWithMsg("用户信息修改成功", c)
}
