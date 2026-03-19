package focus_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
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
		res.SuccessWithMsg("不能关注自己", c)
		return
	}

	//查被关注的人是否存在
	var user models.UserModel
	err := global.Db.Take(&user, cr.FocusUserID).Error
	if err != nil {
		res.FailWithMsg("关注用户不存在", c)
		return
	}

	// 查之前是否已经关注了他
	var fous models.UserFocusModel
	err = global.Db.Unscoped().Take(&fous, "user_id = ? and focus_user_id = ?", userID, cr.FocusUserID).Error
	if err == nil && fous.DeletedAt.Time.IsZero() {
		res.FailWithMsg("已经关注了", c)
		return
	}

	// 关注是否有限制？其做法可用 redis 自增检查，设置过期时间为今天的 23:59:59

	// 关注 - 如果存在且被软删除则恢复，否则创建新记录
	if !fous.DeletedAt.Time.IsZero() {
		// 记录存在但被软删除，恢复它（清除 DeletedAt）
		fous.DeletedAt = gorm.DeletedAt{}
		err = global.Db.Save(&fous).Error
		if err != nil {
			res.FailWithMsg("关注恢复失败", c)
			return
		}
	} else {
		// 记录不存在，创建新记录
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
		res.SuccessWithMsg("不能取关自己", c)
		return
	}

	//查被关注的人是否存在
	var user models.UserModel
	err := global.Db.Take(&user, cr.FocusUserID).Error
	if err != nil {
		res.FailWithMsg("取关用户不存在", c)
		return
	}

	// 查之前是否已经关注了他
	var fous models.UserFocusModel
	err = global.Db.Take(&fous, "user_id = ? and focus_user_id = ?", userID, cr.FocusUserID).Error
	if err != nil {
		res.FailWithMsg("未关注该用户", c)
		return
	}

	// 关注是否有限制？ 其做法可用redis 自增检查 ， 设置过期时间为今天的23：59：59

	// 关注
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

	if cr.UserID != 0 && userID != cr.UserID { // 排除自己情况的限制
		var user models.UserConfigModel
		err := global.Db.Take(&user, "user_id = ? ", cr.UserID).Error
		if err != nil {
			res.FailWithMsg("用户配置信息不存在", c)
			return
		}

		if !user.OpenFollow {
			res.FailWithMsg("用户未开放我的关注", c)
			return
		}
	} else {
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			res.FailWithMsg("请登录", c)
			return
		}
		cr.UserID = claims.Claims.UserID
	}

	_list, count, _ := common.ListQuery(models.UserFocusModel{
		UserID:      cr.UserID,
		FocusUserID: cr.FocusUserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Preloads: []string{"FocusUserModel"}})

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
			res.FailWithMsg("用户配置信息不存在", c)
			return
		}

		if !user.OpenFans {
			res.FailWithMsg("用户未开放我的粉丝", c)
			return
		}
	} else {
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			res.FailWithMsg("请登录", c)
			return
		}
		cr.UserID = claims.Claims.UserID
	}

	_list, count, _ := common.ListQuery(models.UserFocusModel{
		UserID:      cr.FocusUserID,
		FocusUserID: cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Preloads: []string{"FocusUserModel"}})

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
