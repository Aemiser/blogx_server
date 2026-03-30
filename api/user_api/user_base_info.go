package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/redis_service/redis_user"

	"github.com/gin-gonic/gin"
)

type UserBaseInfoResponse struct {
	id           uint   `json:"id"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	CodeAge      uint   `json:"codeAge"`
	LookCount    int    `json:"lookCount"`
	LikeCount    uint   `json:"likeCount"`
	FollowCount  uint   `json:"followCount"`
	FansCount    uint   `json:"fansCount"`
	CollectCount int    `json:"collectCount"`
	ArticleCount int    `json:"articleCount"`
	Place        string `json:"place"`
	OpenCollect  bool   `json:"openCollect"` // 公开我的收藏
	OpenFollow   bool   `json:"openFollow"`  // 公开我的关注
	OpenFans     bool   `json:"openFans"`    // 公开我的粉丝
	HomeStyleID  uint   `json:"homeStyleID"` // 主页样式ID
}

func (UserApi) UserBaseInfoView(c *gin.Context) {
	var req models.IDRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	var userModel models.UserModel
	err = global.Db.Preload("UserConfigModel").Preload("ArticleList").Take(&userModel, req.ID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	data := UserBaseInfoResponse{
		id:        userModel.ID,
		Nickname:  userModel.Nickname,
		Avatar:    userModel.Avatar,
		CodeAge:   userModel.GetCodeAge(),
		LookCount: userModel.UserConfigModel.LookCount + redis_user.GetUserCacheLook(req.ID),
		//LikeCount:    1, //TODO 获取用户点赞数
		FollowCount: 0,
		FansCount:   0,
		//CollectCount: 1, //TODO 获取用户收藏数
		ArticleCount: len(userModel.ArticleList),
		Place:        userModel.Addr,
		OpenCollect:  userModel.UserConfigModel.OpenCollect,
		OpenFollow:   userModel.UserConfigModel.OpenFollow,
		OpenFans:     userModel.UserConfigModel.OpenFans,
		HomeStyleID:  userModel.UserConfigModel.HomeStyleID,
	}

	var fousList []models.UserFocusModel
	global.Db.Find(&fousList, "focus_user_iD = ? or user_id", req.ID, req.ID)
	for _, model := range fousList {
		if model.UserID == req.ID {
			data.FollowCount++
		} else {
			data.FansCount++
		}
	}

	redis_user.SetCacheLook(req.ID, true)
	res.SuccessWithData(data, c)

}
