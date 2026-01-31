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
	StartTime string `json:"startTime" form:"startTime"` // 起止时间的时间戳
	EndTime   string `json:"endTime" form:"endTime"`
	Type      int8   `json:"type" form:"type" binding:"required,oneof=1 2"`
}

type UserLoginListResponse struct {
	models.UserLoginModel
	UserNickname string `json:"userNickname,omitempty"`
	UserAvatar   string `json:"userAvatar,omitempty"`
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

	if req.StartTime > "" {
		_, err = time.Parse("2006-01-02 15:04:05", req.StartTime)
		if err != nil {
			res.FailWithMsg("起始时间格式错误", c)
			return
		}
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime > "" {
		_, err = time.Parse("2006-01-02 15:04:05", req.EndTime)
		if err != nil {
			res.FailWithMsg("截至时间格式错误", c)
			return
		}
		query = query.Where("created_at <= ?", req.EndTime)
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
