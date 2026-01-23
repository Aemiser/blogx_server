package log_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	e "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type runtimeDateType int8

const (
	RuntimedateHour  runtimeDateType = 1
	RuntimedateDay   runtimeDateType = 2
	RuntimedateWeek  runtimeDateType = 3
	RuntimedateMonth runtimeDateType = 4
)

func (r runtimeDateType) GetSqlTime() string {
	switch r {
	case RuntimedateHour:
		return "interval 1 HOUR"
	case RuntimedateDay:
		return "interval 1 DAY"
	case RuntimedateWeek:
		return "interval 1 WEEK"
	case RuntimedateMonth:
		return "interval 1 MONTH"
	}
	return "interval 1 DAY"
}

type RuntimeLog struct {
	title           string
	level           enum.LogLevelType
	serviceName     string
	itemList        []string
	runtimeDateType runtimeDateType
}

func NewRuntimeLog(serviceName string, runtimeType runtimeDateType) *RuntimeLog {
	return &RuntimeLog{
		serviceName:     serviceName,
		runtimeDateType: runtimeType,
	}
}
func (ac *RuntimeLog) Save() {

	// 判断是创建还是更新
	var log models.LogModel
	global.Db.Find(&log, fmt.Sprintf("service_name = ? and log_type = ?  and created_at >= date_sub(now(),%s)", ac.runtimeDateType.GetSqlTime()),
		ac.serviceName,
		enum.RuntimeLogType)

	content := strings.Join(ac.itemList, "\n")
	// 存在即更新
	if log.ID != 0 {
		newContent := log.Content + "\n" + content

		global.Db.Model(&log).Updates(map[string]any{
			"content": newContent,
		})
		ac.itemList = []string{}
		return
	}

	// 不存在即创建
	err := global.Db.Create(&models.LogModel{
		LogType:     enum.RuntimeLogType,
		Content:     content,
		Title:       ac.title,
		Level:       ac.level,
		ServiceName: ac.serviceName,
	}).Error

	if err != nil {
		logrus.Errorf("创建运行日志失败：%s", err.Error())
	}
	ac.itemList = []string{}

	return
}

func (ac *RuntimeLog) setItem(label string, value any, loglevel enum.LogLevelType) {
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

func (ac *RuntimeLog) SetItem(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}

func (ac *RuntimeLog) SetItemInfo(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}

func (ac *RuntimeLog) SetItemWarn(label string, value any) {
	ac.setItem(label, value, enum.LogWarnLevel)
}
func (ac *RuntimeLog) SetItemError(label string, value any) {
	ac.setItem(label, value, enum.LogErrLevel)
}

func (ac *RuntimeLog) SetError(label string, err error) {
	msg := e.WithStack(err)
	logrus.Errorf(err.Error())
	ac.itemList = append(ac.itemList, fmt.Sprintf("\n<div class=\"log_error\">\n    <div class=\"line\">\n        <div class=\"label\">%s</div>\n        <div class=\"value\">%s</div>\n        <div class=\"type\">%T</div>\n    </div>\n    <div class=\"stack\">%+v</div>\n</div>\n",
		label, // 错误信息
		err,   // 错误内容
		err,   // 错误类型
		msg,   // 错误堆栈
	))
}
