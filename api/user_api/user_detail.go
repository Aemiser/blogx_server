package user_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"fmt"
	"math"
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
	CodeAge        int                     `json:"codeAge"`        //码龄
	models.UserConfigModel
}

func (UserApi) UserDetailView(c *gin.Context) {
	claims := jwts.GetClaimsByGin(c)
	fmt.Println("claims:", claims)
	var userModel models.UserModel
	err := global.Db.Preload("UserConfigModel").Take(&userModel, claims.Claims.UserID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}
	// 计算码龄
	sub := time.Now().Sub(userModel.CreatedAt)
	CodeAge := int(math.Ceil(sub.Hours() / 24 / 365))

	var data = UserDetailResponse{
		ID:             userModel.ID,
		CreatedAt:      userModel.CreatedAt,
		Username:       userModel.Username,
		Nickname:       userModel.Nickname,
		Avatar:         userModel.Avatar,
		Abstract:       userModel.Abstract,
		Email:          userModel.Email,
		RegisterSource: userModel.RegisterSource,
		CodeAge:        CodeAge,
	}

	if userModel.UserConfigModel != nil {
		data.UserConfigModel = *userModel.UserConfigModel
	}

	res.SuccessWithData(data, c)
}
