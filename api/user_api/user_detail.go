package user_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"time"

	"github.com/gin-gonic/gin"
)

type UserDetailResponse struct {
	ID             uint                    `json:"id"`
	CreatedAt      time.Time               `json:"createdAt"`
	Username       string                  `json:"username"`
	Nickname       string                  `json:"nickname"`
	Avatar         string                  `json:"avatar"`
	Abstract       string                  `json:"abstract"`
	Email          string                  `gorm:"size:256" json:"email"`
	RegisterSource enum.RegisterSourceType `json:"registerSource"` //注册来源
	CodeAge        uint                    `json:"codeAge"`        //码龄
	Role           enum.RoleType           `json:"role"`
	models.UserConfigModel
	UserPassword bool `json:"userPassword"`
}

func (UserApi) UserDetailView(c *gin.Context) {
	claims := jwts.GetClaimsByGin(c)
	var userModel models.UserModel
	err := global.Db.Preload("UserConfigModel").Take(&userModel, claims.Claims.UserID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	var data = UserDetailResponse{
		ID:             userModel.ID,
		CreatedAt:      userModel.CreatedAt,
		Username:       userModel.Username,
		Nickname:       userModel.Nickname,
		Avatar:         userModel.Avatar,
		Abstract:       userModel.Abstract,
		Email:          userModel.Email,
		RegisterSource: userModel.RegisterSource,
		CodeAge:        userModel.GetCodeAge(),
		Role:           userModel.Role,
	}
	if userModel.Password != "" {
		data.UserPassword = true
	}

	if userModel.UserConfigModel != nil {
		data.UserConfigModel = *userModel.UserConfigModel
	}

	res.SuccessWithData(data, c)
}
