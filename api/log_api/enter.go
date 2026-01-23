package log_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/models"
	"blogx_server/models/enum"
	"fmt"

	"github.com/gin-gonic/gin"
)

type LogApi struct {
}

type LogListRequest struct {
	common.PageInfo
	LogType     enum.LogType      `form:"logType" json:"logType"`
	Level       enum.LogLevelType `form:"level" json:"level"`
	UserID      uint              `form:"userID" json:"userID"`
	IP          string            `form:"ip" json:"ip"`
	LoginStatus bool              `form:"loginStatus" json:"loginStatus"`
	ServiceName string            `form:"serviceName" json:"serviceName"`
}

type LogListResponse struct {
	models.LogModel
	UserNickName string `json:"userNickName"`
	UserAvatar   string `json:"userAvatar"`
}

func (LogApi) LogListView(c *gin.Context) {
	// 分页查询 精确查询，模糊匹配
	var req LogListRequest
	err := c.ShouldBind(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	fmt.Println(req)
	list, count, err := common.ListQuery(models.LogModel{
		LogType:     req.LogType,
		Level:       req.Level,
		UserID:      req.UserID,
		IP:          req.IP,
		LoginStatus: req.LoginStatus,
		ServiceName: req.ServiceName,
	}, common.Options{

		PageInfo:     req.PageInfo,
		Preloads:     []string{"UserModel"},
		Likes:        []string{"Title"},
		Debug:        true,
		DefaultOrder: "created_at desc",
	},
	)
	//
	//var List []models.LogModel
	//if req.Page >= 20 {
	//	req.Page = 1
	//}
	//if req.Page <= 0 {
	//	req.Page = 1
	//}
	//
	//if req.Limit == 0 || req.Limit > 100 {
	//	req.Limit = 10
	//}
	//offest := (req.Page - 1) * req.Limit
	//model := models.LogModel{
	//	LogType:     req.LogType,
	//	Level:       req.Level,
	//	UserID:      req.UserID,
	//	IP:          req.IP,
	//	LoginStatus: req.LoginStatus,
	//	ServiceName: req.ServiceName,
	//}
	//
	//like := global.Db.Debug().Where("title like ?", fmt.Sprintf("%%%s%%", req.Key))
	//global.Db.Debug().Preload("UserModel").Where(like).Where(model).Offset(offest).Limit(req.Limit).Find(&List)
	//
	//var count int64
	//global.Db.Debug().Where(like).Where(model).Model(&models.LogModel{}).Count(&count)

	var _list = make([]LogListResponse, 0)
	for _, logModel := range list {
		_list = append(_list, LogListResponse{
			LogModel:     logModel,
			UserNickName: logModel.UserModel.Nickname,
			UserAvatar:   logModel.UserModel.Avatar,
		})

	}
	res.FailWithList(_list, int(count), c)
	return
}
