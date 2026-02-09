package user_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/maps"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type UserInfoUpdataRequest struct {
	Username    *string   `json:"username" s-u:"username"`
	Nickname    *string   `json:"nickname" s-u:"nickname"`
	Avatar      *string   `json:"avatar" s-u:"avatar"`
	Abstract    *string   `json:"abstract" s-u:"abstract"`
	LikeTags    *[]string `json:"likeTags" s-u-c:"like_tags"`
	OpenCollect *bool     `json:"openCollect" s-u-c:"open_collect"`  // 公开我的收藏
	OpenFollow  *bool     `json:"openFollow" s-u-c:"open_follow"`    // 公开我的关注
	OpenFans    *bool     `json:"openFans" s-u-c:"open_fans"`        // 公开我的粉丝
	HomeStyleID *uint     `json:"homeStyleID" s-u-c:"home_style_id"` // 主页样式ID
}

func (UserApi) UserInfoUpdate(c *gin.Context) {
	var req UserInfoUpdataRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	userMap := maps.StructToMap(req, "s-u")
	userConfMap := maps.StructToMap(req, "s-u-c")

	fmt.Println(userMap)
	fmt.Println(userConfMap)

	claims := jwts.GetClaimsByGin(c)
	if len(userMap) > 0 {
		var userModel models.UserModel
		err = global.Db.Preload("UserConfigModel").Take(&userModel, claims.Claims.UserID).Error
		if err != nil {
			res.FailWithMsg("用户不存在", c)
			return
		}

		// 判断
		if req.Username != nil {
			var userCount int64
			// 判断用户名是否被其他人使用
			global.Db.Model(&models.UserModel{}).
				Where("username = ? and id <> ? ", req.Username, claims.Claims.UserID).
				Count(&userCount)
			if userCount > 0 {
				res.FailWithMsg("该用户名已被使用", c)
				return
			}
			// 用户名相同
			if *req.Username != userModel.Username {
				// 判断时间期限
				var uud = userModel.UserConfigModel.UpdataUsernameDate
				if uud != nil {
					if time.Now().Sub(*uud).Hours() < 24*30 {
						res.FailWithMsg("用户名30天内只能修改一次", c)
						return
					}
				}
			}
			// 大于30 天后可以更改，修改updata字段
			userConfMap["updata_username_date"] = time.Now()
		}

		if req.Nickname != nil || req.Avatar != nil {
			if userModel.RegisterSource == enum.RegisterSourceTypeQQ {
				res.FailWithMsg("QQ用户不允许修改昵称和头像", c)
				return
			}
		}

		err = global.Db.Model(&userModel).Updates(userMap).Error
		if err != nil {
			res.FailWithMsg("用户信息修改失败", c)
			return
		}
	}

	if len(userConfMap) > 0 {
		// 查询用户
		var userConfModel models.UserConfigModel
		err = global.Db.Take(&userConfModel, "user_id =?", claims.Claims.UserID).Error
		if err != nil {
			res.FailWithMsg("用户不存在", c)
			return
		}
		// 修改用户
		err = global.Db.Model(&userConfModel).Updates(userConfMap).Error
		if err != nil {
			res.FailWithMsg("用户信息修改失败", c)
			return
		}
	}

	res.SuccessWithMsg("用户信息修改成功", c)
}
