package res

import (
	"blogx_server/utils/validata"

	"github.com/gin-gonic/gin"
)

type Code int

const (
	SuccessCode     Code = 0
	FailValueCode   Code = 1001
	FailServiceCode Code = 1002
)

func (c Code) String() string {
	switch c {
	case SuccessCode:
		return "成功"
	case FailValueCode:
		return "参数错误"
	case FailServiceCode:
		return "服务错误"
	default:
		return "未知错误"
	}
}

type Response struct {
	Code Code   `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

var empty = map[string]any{}

func (r Response) Json(c *gin.Context) {
	c.JSON(200, r)
}
func Success(data any, msg string, c *gin.Context) {
	Response{SuccessCode, data, msg}.Json(c)
}

func SuccessWithData(data any, c *gin.Context) {
	Response{SuccessCode, data, "Success"}.Json(c)
}

func SuccessWithMsg(msg string, c *gin.Context) {
	Response{SuccessCode, empty, msg}.Json(c)
}

func FailWithList(list any, count int, c *gin.Context) {
	Response{FailValueCode, map[string]any{
		"list":  list,
		"count": count,
	}, "Success"}.Json(c)
}

func FailWithData(data any, msg string, c *gin.Context) {
	Response{FailValueCode, data, msg}.Json(c)
}

func FailWithMsg(msg string, c *gin.Context) {
	Response{FailValueCode, empty, msg}.Json(c)
}

func FailWithCode(code Code, c *gin.Context) {
	Response{code, empty, code.String()}.Json(c)
}

func FailWithError(err error, c *gin.Context) {
	data, msg := validata.ValidateError(err)
	FailWithData(data, msg, c)
}
