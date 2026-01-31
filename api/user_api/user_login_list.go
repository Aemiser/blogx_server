package user_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"time"

	"github.com/gin-gonic/gin"
)

type UserLoginListRequest struct {
	common.PageInfo
	UserID    uint   `json:"userID" form:"userID"`
	Ip        string `json:"ip" form:"ip"`
	Addr      string `json:"addr" form:"addr"`
	StartTime int64  `json:"startTime" form:"startTime"` // 起止时间的时间戳
	EndTime   int64  `json:"endTime" form:"endTime"`
	Type      int8   `json:"type" form:"type" binding:"required,oneof=1 2"`
}

type UserLoginListResponse struct {
	models.UserLoginModel
	UserNickname string `json:"userNickname"`
	UserAvatar   string `json:"userAvatar"`
}

func (UserApi) UserLoginListView(c *gin.Context) {
	var req UserLoginListRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	claims := jwts.GetClaimsByGin(c)
	if req.Type == 1 {
		req.UserID = claims.Claims.UserID
	}

	var query = global.Db.Where("")

	if req.StartTime > 0 {
		t := time.Unix(req.StartTime, 0)
		query = query.Where("created_at >= ?", t)
	}
	if req.EndTime > 0 {
		t := time.Unix(req.EndTime, 0)
		query = query.Where("created_at <= ?", t)
	}
	var preloads []string
	if req.Type == 2 {
		preloads = []string{"UserModel"}
	}

	_list, count, _ := common.ListQuery[models.UserLoginModel](models.UserLoginModel{
		UserID: req.UserID,
		IP:     req.Ip,
		Addr:   req.Addr,
	}, common.Options{
		PageInfo: req.PageInfo,
		Preloads: preloads,
		Where:    query,
	})

	var list = make([]UserLoginListResponse, 0)
	for _, model := range _list {
		list = append(list, UserLoginListResponse{
			UserLoginModel: model,
			UserNickname:   model.UserModel.Nickname,
			UserAvatar:     model.UserModel.Avatar,
		})
	}

	res.SuccessWithList(list, count, c)

}
