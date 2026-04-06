package log_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/log_service"
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
	})

	var _list = make([]LogListResponse, 0)
	for _, logModel := range list {
		_list = append(_list, LogListResponse{
			LogModel:     logModel,
			UserNickName: logModel.UserModel.Nickname,
			UserAvatar:   logModel.UserModel.Avatar,
		})

	}
	res.SuccessWithList(_list, int(count), c)
	return
}

func (LogApi) LogReadView(c *gin.Context) {
	var req models.IDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	var log models.LogModel
	err := global.Db.Take(&log, req.ID).Error
	if err != nil {
		res.FailWithMsg("不存在的日志", c)
		return
	}

	// 如果未读则修改
	if !log.IsRead {
		global.Db.Model(&log).Update("is_read", true)
	}

	res.SuccessWithMsg("日志读取成功", c)
	return
}

func (LogApi) LogRemoveView(c *gin.Context) {
	var req models.IDListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	log := log_service.GetLog(c)
	log.ShowResponse()
	log.ShowRequest()
	var logList []models.LogModel
	global.Db.Find(&logList, "id in ?", req.IDList)

	if len(logList) > 0 {
		global.Db.Delete(&logList)
	}
	msg := fmt.Sprintf("共删除%d条日志", len(logList))
	res.SuccessWithMsg(msg, c)

}
