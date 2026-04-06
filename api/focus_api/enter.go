package focus_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FocusApi struct {
}

type FocusUserRequest struct {
	FocusUserID uint `json:"focusUserID" binding:"required"`
}

// FocusUserApi 登录人关注用户
func (FocusApi) FocusUserApi(c *gin.Context) {
	cr := middlerware.GetBind[FocusUserRequest](c)
	userID := jwts.GetUserIDByGin(c)
	if userID == cr.FocusUserID {
		res.FailWithCodeAndMsg(res.FocusCannotSelf, "不能关注自己", c)
		return
	}

	var user models.UserModel
	err := global.Db.Take(&user, cr.FocusUserID).Error
	if err != nil {
		res.FailWithCode(res.FocusUserNotFound, c)
		return
	}

	var fous models.UserFocusModel
	err = global.Db.Unscoped().Take(&fous, "user_id = ? and focus_user_id = ?", userID, cr.FocusUserID).Error
	if err == nil && fous.DeletedAt.Time.IsZero() {
		res.FailWithCode(res.FocusAlreadyExists, c)
		return
	}

	if !fous.DeletedAt.Time.IsZero() {
		fous.DeletedAt = gorm.DeletedAt{}
		err = global.Db.Save(&fous).Error
		if err != nil {
			res.FailWithCode(res.FocusRestoreFailed, c)
			return
		}
	} else {
		global.Db.Create(&models.UserFocusModel{
			UserID:      userID,
			FocusUserID: cr.FocusUserID,
		})
	}

	res.SuccessWithMsg("关注成功", c)
	return
}

func (FocusApi) UnFocusUserApi(c *gin.Context) {
	cr := middlerware.GetBind[FocusUserRequest](c)
	userID := jwts.GetUserIDByGin(c)
	if userID == cr.FocusUserID {
		res.FailWithCodeAndMsg(res.FocusCannotSelfUn, "不能取关自己", c)
		return
	}

	var user models.UserModel
	err := global.Db.Take(&user, cr.FocusUserID).Error
	if err != nil {
		res.FailWithCode(res.UnFocusUserNotFound, c)
		return
	}

	var fous models.UserFocusModel
	err = global.Db.Take(&fous, "user_id = ? and focus_user_id = ?", userID, cr.FocusUserID).Error
	if err != nil {
		res.FailWithCode(res.UnFocusNotFollowed, c)
		return
	}

	global.Db.Delete(&fous)

	res.SuccessWithMsg("取关成功", c)
	return
}

type FocusUserListRequest struct {
	common.PageInfo
	FocusUserID uint `form:"focusUserID" `
	UserID      uint `form:"userID"`
}

type FocusUserListResponse struct {
	FocusUserID       uint      `json:"focusUserID"`
	FocusUserNickname string    `json:"focusUserNickname"`
	FocusUserAvatar   string    `json:"focusUserAvatar"`
	FocusUserAbstract string    `json:"focusUserAbstract"`
	CreateAt          time.Time `json:"createAt"`
}

func (FocusApi) FocusUserListApi(c *gin.Context) {
	cr := middlerware.GetBind[FocusUserListRequest](c)
	userID := jwts.GetUserIDByGin(c)
	claims, err := jwts.ParseTokenByGin(c)
	if cr.UserID != 0 && userID != cr.UserID {
		var user models.UserConfigModel
		err1 := global.Db.Take(&user, "user_id = ? ", cr.UserID).Error
		if err1 != nil {
			res.FailWithCode(res.FocusUserConfigNot, c)
			return
		}

		if !user.OpenFollow {
			res.FailWithCode(res.FocusUserNotOpen, c)
			return
		}

		if claims == nil && err != nil {
			if cr.Limit > 10 || cr.Page > 1 {
				res.FailWithCode(res.SysUnauthorized, c)
				return
			}

		}
	} else {
		if err != nil {
			res.FailWithCode(res.SysUnauthorized, c)
			return
		}
		cr.UserID = claims.Claims.UserID
	}

	query := global.Db.Where("")
	if cr.Key != "" {
		var userIDList []uint
		global.Db.Model(models.UserModel{}).Where("nickname like ?", fmt.Sprintf("%%%s%%", cr.Key)).
			Select("id").Scan(&userIDList)

		if len(userIDList) > 0 {
			query.Where("focus_user_id in ?", userIDList)
		}
	}

	_list, count, _ := common.ListQuery(models.UserFocusModel{
		UserID:      cr.UserID,
		FocusUserID: cr.FocusUserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Preloads: []string{"FocusUserModel"},
		Where:    query,
	})

	var list = make([]FocusUserListResponse, 0)
	for _, model := range _list {
		list = append(list, FocusUserListResponse{
			FocusUserID:       model.FocusUserID,
			FocusUserNickname: model.FocusUserModel.Nickname,
			FocusUserAvatar:   model.FocusUserModel.Avatar,
			FocusUserAbstract: model.FocusUserModel.Abstract,
			CreateAt:          model.CreatedAt,
		})
	}
	res.SuccessWithList(list, count, c)
}

type FansUserListResponse struct {
	FansUserID       uint      `json:"fansUserID"`
	FansUserNickname string    `json:"fansUserNickname"`
	FansUserAvatar   string    `json:"fansUserAvatar"`
	FansUserAbstract string    `json:"fansUserAbstract"`
	CreateAt         time.Time `json:"createAt"`
}

func (FocusApi) FansUserListApi(c *gin.Context) {
	cr := middlerware.GetBind[FocusUserListRequest](c)

	if cr.UserID != 0 {
		var user models.UserConfigModel
		err := global.Db.Take(&user, "user_id = ? ", cr.UserID).Error
		if err != nil {
			res.FailWithCode(res.FansUserConfigNot, c)
			return
		}

		if !user.OpenFans {
			res.FailWithCode(res.FansUserNotOpen, c)
			return
		}
	} else {
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			res.FailWithCode(res.SysUnauthorized, c)
			return
		}
		cr.UserID = claims.Claims.UserID
	}

	_list, count, _ := common.ListQuery(models.UserFocusModel{
		UserID:      cr.FocusUserID,
		FocusUserID: cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Preloads: []string{"UserModel"}})

	var list = make([]FansUserListResponse, 0)
	for _, model := range _list {
		list = append(list, FansUserListResponse{
			FansUserID:       model.UserID,
			FansUserNickname: model.UserModel.Nickname,
			FansUserAvatar:   model.UserModel.Avatar,
			FansUserAbstract: model.UserModel.Abstract,
			CreateAt:         model.CreatedAt,
		})
	}
	res.SuccessWithList(list, count, c)
}
