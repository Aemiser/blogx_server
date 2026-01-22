package log_service

import (
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
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

func (ac *ActionLog) setItem(label string, value any, loglevel enum.LogLevelType) {
	var v string
	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		byteDate, _ := json.Marshal(value)
		v = string(byteDate)
	default:
		v = fmt.Sprintf("%v", value)
	}
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_item %s\">\n    <div class=\"log_item_label\">%s</div>\n    <div class=\"log_item_content\">%s</div>\n</div>",
		loglevel,
		label,
		v,
	))
}

func (ac *ActionLog) SetItem(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}

func (ac *ActionLog) SetItemInfo(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}

func (ac *ActionLog) SetItemWarn(label string, value any) {
	ac.setItem(label, value, enum.LogWarnLevel)
}
func (ac *ActionLog) SetItemError(label string, value any) {
	ac.setItem(label, value, enum.LogErrLevel)
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

	tmpItemList := []string{}
	// 设置请求
	if ac.showRequest {
		tmpItemList = append(tmpItemList, fmt.Sprintf("<div class=\"log_request\">\n    <div class=\"log_request_head\">\n        <span class=\"log_request_method delete\">%s</span>\n        <span class=\"log_request_path\">%s</span>\n    </div>\n    <div class=\"log_request_body\">\n        <pre class=\"log_json_body\">%s</pre>\n    </div>\n</div>",
			ac.c.Request.Method,
			ac.c.Request.URL.String(),
			string(ac.RequestBody),
		))
	}

	// 中间contest
	tmpItemList = append(tmpItemList, ac.itemList...)

	// 设置响应
	if ac.showResponse {
		tmpItemList = append(tmpItemList, fmt.Sprintf("<div class=\"log_response\">\n    <pre class=\"log_json_body\">%s</pre>\n</div>",
			string(ac.ResponseBody),
		))
	}

	log := models.LogModel{
		LogType: enum.ActionLogType,
		Title:   ac.title,
		Content: strings.Join(tmpItemList, "\n"),
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
