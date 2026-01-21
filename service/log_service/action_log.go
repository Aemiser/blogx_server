package log_service

import (
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ActionLog struct {
	c            *gin.Context
	title        string
	level        enum.LogLevelType
	RequestBody  []byte
	ResponseBody []byte
	log          *models.LogModel
}

func NewActionLog(c *gin.Context) *ActionLog {
	return &ActionLog{c: c}
}

func GetLog(c *gin.Context) *ActionLog {
	_log, ok := c.Get("log")
	if !ok {
		return NewActionLog(c)
	}
	log, ok := _log.(*ActionLog)
	if !ok {
		return NewActionLog(c)
	}
	return log

}

func (ac *ActionLog) SetTitle(title string) {
	ac.title = title
}

func (ac *ActionLog) SetLevel(level enum.LogLevelType) {
	ac.level = level
}

func (ac *ActionLog) SetRequest(c *gin.Context) {
	byteData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(byteData))
	ac.RequestBody = byteData
}

func (ac *ActionLog) SetResponse(data []byte) {
	ac.ResponseBody = data
}

func (ac *ActionLog) Save() {
	if ac.log != nil {
		// 之前创建了，下次就是更新
		global.Db.Model(ac.log).Updates(map[string]any{
			"title": "更新",
		})
		return
	}

	ip := ac.c.ClientIP()
	addr := core.GetIpAddr(ip)
	userID := uint(1)

	log := models.LogModel{
		LogType: enum.ActionLogType,
		Title:   ac.title,
		Content: "",
		Level:   ac.level,
		UserID:  userID,
		IP:      ip,
		Addr:    addr,
	}
	err := global.Db.Create(&log).Error
	if err != nil {
		logrus.Errorf("日志创建失败 %s", err)
	}
	ac.log = &log

}
