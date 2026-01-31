package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

type UserBaseInfoResponse struct {
	id           uint   `json:"id"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	CodeAge      uint   `json:"codeAge"`
	LookCount    uint   `json:"lookCount"`
	LikeCount    uint   `json:"likeCount"`
	FollowCount  uint   `json:"followCount"`
	FansCount    uint   `json:"fansCount"`
	CollectCount uint   `json:"collectCount"`
	ArticleCount uint   `json:"articleCount"`
	Place        string `json:"place"`
}

func (UserApi) UserBaseInfoView(c *gin.Context) {
	var req models.IDRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	var userModel models.UserModel
	err = global.Db.Take(&userModel, req.ID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	data := UserBaseInfoResponse{
		id:           userModel.ID,
		Nickname:     userModel.Nickname,
		Avatar:       userModel.Avatar,
		CodeAge:      userModel.GetCodeAge(),
		LookCount:    1, // TODO:做完文章浏览关系回来写
		LikeCount:    1,
		FollowCount:  1, // TODO:做完好友关系回来写
		FansCount:    1,
		CollectCount: 1,
		ArticleCount: 1,
		Place:        userModel.Addr,
	}
	res.SuccessWithData(data, c)

}
