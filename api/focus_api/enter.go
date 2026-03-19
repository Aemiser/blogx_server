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
	err = global.Db.Take(&fous, "user_id = ? and focus_user_id = ?", userID, cr.FocusUserID).Error
	if err == nil {
		res.FailWithMsg("已经关注了", c)
		return
	}

	// 关注是否有限制？ 其做法可用redis 自增检查 ， 设置过期时间为今天的23：59：59

	// 关注
	global.Db.Create(&models.UserFocusModel{
		UserID:      userID,
		FocusUserID: cr.FocusUserID,
	})

	res.SuccessWithMsg("关注成功", c)
	return
}

type FocusUserListRequest struct {
	common.PageInfo
	FocusUserID uint `json:"focusUserID" `
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
	_list, count, _ := common.ListQuery(models.UserFocusModel{
		UserID:      userID,
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
