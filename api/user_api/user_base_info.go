package user_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum/relationship_enum"
	"blogx_server/service/focus_service"
	"blogx_server/service/redis_service/redis_user"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserBaseInfoResponse struct {
	ID           uint                       `json:"id"`
	Nickname     string                     `json:"nickname"`
	Avatar       string                     `json:"avatar"`
	CodeAge      uint                       `json:"codeAge"`
	LookCount    int                        `json:"lookCount"`
	LikeCount    uint                       `json:"likeCount"`
	FollowCount  uint                       `json:"followCount"`
	FansCount    uint                       `json:"fansCount"`
	CollectCount int                        `json:"collectCount"`
	ArticleCount int                        `json:"articleCount"`
	Place        string                     `json:"place"`
	OpenCollect  bool                       `json:"openCollect"` // 公开我的收藏
	OpenFollow   bool                       `json:"openFollow"`  // 公开我的关注
	OpenFans     bool                       `json:"openFans"`    // 公开我的粉丝
	HomeStyleID  uint                       `json:"homeStyleID"` // 主页样式 ID
	Relation     relationship_enum.Relation `json:"relation"`
}

func (UserApi) UserBaseInfoView(c *gin.Context) {
	var req models.IDRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	var userModel models.UserModel
	err = global.Db.Debug().Preload("ArticleList").Take(&userModel, req.ID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	// 检查 UserConfigModel 是否存在
	var lookCount int
	var openCollect, openFans, openFollow bool
	var homeStyleID uint
	if userModel.UserConfigModel != nil {
		lookCount = userModel.UserConfigModel.LookCount
		openCollect = userModel.UserConfigModel.OpenCollect
		openFollow = userModel.UserConfigModel.OpenFollow
		openFans = userModel.UserConfigModel.OpenFans
		homeStyleID = userModel.UserConfigModel.HomeStyleID
	}

	data := UserBaseInfoResponse{
		ID:           userModel.ID,
		Nickname:     userModel.Nickname,
		Avatar:       userModel.Avatar,
		CodeAge:      userModel.GetCodeAge(),
		LookCount:    lookCount + redis_user.GetUserCacheLook(req.ID),
		FollowCount:  0,
		FansCount:    0,
		ArticleCount: len(userModel.ArticleList),
		Place:        userModel.Addr,
		OpenCollect:  openCollect,
		OpenFollow:   openFollow,
		OpenFans:     openFans,
		HomeStyleID:  homeStyleID,
	}

	// 用户主页的关注信息
	claims, err := jwts.ParseTokenByGin(c)
	if err == nil && claims != nil {
		data.Relation = focus_service.CalcUserRelationship(claims.Claims.UserID, req.ID)
		fmt.Println("用户关系", data.Relation)
	}

	var fousList []models.UserFocusModel
	global.Db.Find(&fousList, "focus_user_iD = ? or user_id = ?", req.ID, req.ID)
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
