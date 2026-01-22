package log_service

import (
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"bytes"
	"fmt"
	"io"
	"strings"

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
	showRequest  bool
	showResponse bool
	itemList     []string
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

func (ac *ActionLog) ShowResponse() {
	ac.showResponse = true
}

func (ac *ActionLog) ShowRequest() {
	ac.showRequest = true
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

	// 设置请求
	if ac.showRequest {
		ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_request\">\n    <div class=\"log_request_head\">\n        <span class=\"log_request_method delete\">%s</span>\n        <span class=\"log_request_path\">%s</span>\n    </div>\n    <div class=\"log_request_body\">\n        <pre class=\"log_json_body\">%s</pre>\n    </div>\n</div>",
			ac.c.Request.Method,
			ac.c.Request.URL.String(),
			string(ac.RequestBody),
		))
	}

	// 设置响应
	if ac.showResponse {
		ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_response\">\n    <pre class=\"log_json_body\">%s</pre>\n</div>",
			string(ac.ResponseBody),
		))
	}

	log := models.LogModel{
		LogType: enum.ActionLogType,
		Title:   ac.title,
		Content: strings.Join(ac.itemList, "\n"),
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
