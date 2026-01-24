package log_service

import (
	"blogx_server/common/jwts"
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	e "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ActionLog struct {
	c                  *gin.Context
	title              string
	level              enum.LogLevelType
	RequestBody        []byte
	ResponseBody       []byte
	showRequest        bool
	showResponse       bool
	ResponseHeader     http.Header
	showRequestHeader  bool
	showResponseHeader bool
	log                *models.LogModel
	itemList           []string
	idMiddleware       bool
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
	c.Set("SaveLog", true)
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

func (ac *ActionLog) SetLink(label string, href string) {
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_item link\">\n    <div class=\"log_item_label\">%s</div>\n    <div class=\"log_item_content\">\n        <a href=\"%s\" target=\"_blank\">%s</a>\n    </div>\n</div>\n",
		label,
		href,
		href,
	))

}

func (ac *ActionLog) SetImage(src string) {
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_image\">\n    <img src=\"%s\" alt=\"\">\n</div>",
		src,
	))
}

func (ac *ActionLog) ShowRequestHeader() {
	ac.showRequestHeader = true
}
func (ac *ActionLog) ShowResponseHeader() {
	ac.showResponseHeader = true
}

func (ac *ActionLog) SetResponseHeader(header http.Header) {
	ac.ResponseHeader = header
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

func (ac *ActionLog) SetError(label string, err error) {
	msg := e.WithStack(err)
	logrus.Errorf(err.Error())
	ac.itemList = append(ac.itemList, fmt.Sprintf("\n<div class=\"log_error\">\n    <div class=\"line\">\n        <div class=\"label\">%s</div>\n        <div class=\"value\">%s</div>\n        <div class=\"type\">%T</div>\n    </div>\n    <div class=\"stack\">%+v</div>\n</div>\n",
		label, // 错误信息
		err,   // 错误内容
		err,   // 错误类型
		msg,   // 错误堆栈
	))
}

func (ac *ActionLog) MiddlewareSave() {
	// 没有调用这个对象就不创建日志
	_savelog, _ := ac.c.Get("SaveLog")
	savelog, _ := _savelog.(bool)
	if !savelog {
		return
	}
	if ac.log == nil {
		// 创建
		ac.idMiddleware = true
		ac.Save()
		return
	}
	// 在视图里面Save过，现在更新
	// 设置响应头
	if ac.showResponseHeader {
		byteDate, _ := json.Marshal(ac.ResponseHeader)
		ac.itemList = append(ac.itemList, fmt.Sprintf("\n<div class=\"log_response_header\">\n    <div class=\"log_response_body\">\n        <pre class=\"log_json_body\">%s</pre>\n    </div>\n</div>",
			string(byteDate),
		))
	}
	// 设置响应
	if ac.showResponse {
		ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_response\">\n    <pre class=\"log_json_body\">%s</pre>\n</div>",
			string(ac.ResponseBody),
		))
	}
	ac.Save()
	return
}
func (ac *ActionLog) Save() uint {
	// 优化Save:
	// 方案1：save只能在日志的响应中间件里面调用
	// 方案2：在视图里面调用Save，需要返回日志的ID

	if ac.log != nil {
		newList := strings.Join(ac.itemList, "\n")
		content := ac.log.Content + "\n" + newList
		// 之前创建了，下次就是更新
		global.Db.Model(ac.log).Updates(map[string]any{
			"content": content,
		})
		ac.itemList = []string{}
		return ac.log.ID
	}

	ip := ac.c.ClientIP()
	addr := core.GetIpAddr(ip)

	claim, err := jwts.ParseTokenByGin(ac.c)
	userID := uint(0)
	if err == nil && claim != nil {
		userID = claim.Claims.UserID
	}

	tmpItemList := []string{}

	// 设置请求头
	if ac.showRequestHeader {
		byteDate, _ := json.Marshal(ac.c.Request.Header)
		tmpItemList = append(tmpItemList, fmt.Sprintf("\n<div class=\"log_request_header\">\n    <div class=\"log_request_body\">\n        <pre class=\"log_json_body\">%s</pre>\n    </div>\n</div>",
			string(byteDate),
		))
	}
	// 设置请求
	if ac.showRequest {
		tmpItemList = append(tmpItemList, fmt.Sprintf("<div class=\"log_request\">\n    <div class=\"log_request_head\">\n        <span class=\"log_request_method %s\">%s</span>\n        <span class=\"log_request_path\">%s</span>\n    </div>\n    <div class=\"log_request_body\">\n        <pre class=\"log_json_body\">%s</pre>\n    </div>\n</div>",
			strings.ToLower(ac.c.Request.Method),
			ac.c.Request.Method,
			ac.c.Request.URL.String(),
			string(ac.RequestBody),
		))
	}

	// 中间contest
	tmpItemList = append(tmpItemList, ac.itemList...)

	if ac.idMiddleware {
		// 设置响应头
		if ac.showResponseHeader {
			byteDate, _ := json.Marshal(ac.ResponseHeader)
			tmpItemList = append(tmpItemList, fmt.Sprintf("\n<div class=\"log_response_header\">\n    <div class=\"log_response_body\">\n        <pre class=\"log_json_body\">%s</pre>\n    </div>\n</div>",
				string(byteDate),
			))
		}
		// 设置响应
		if ac.showResponse {
			tmpItemList = append(tmpItemList, fmt.Sprintf("<div class=\"log_response\">\n    <pre class=\"log_json_body\">%s</pre>\n</div>",
				string(ac.ResponseBody),
			))
		}
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
	err = global.Db.Create(&log).Error
	if err != nil {
		logrus.Errorf("日志创建失败 %s", err)
	}
	ac.log = &log            // 保存日志副本
	ac.itemList = []string{} // 清空日志信息
	return ac.log.ID         // 返回日志ID
}
