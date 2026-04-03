package res

import (
	"blogx_server/utils/validata"
	"fmt"

	"github.com/gin-gonic/gin"
)

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

func SuccessWithMsgf(c *gin.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	Response{SuccessCode, empty, msg}.Json(c)
}

func SuccessWithList(list any, count int, c *gin.Context) {
	Response{SuccessCode, map[string]any{
		"list":  list,
		"count": count,
	}, "Success"}.Json(c)
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

func FailWithMsgf(c *gin.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	Response{FailValueCode, empty, msg}.Json(c)
}

func FailWithCode(code Code, c *gin.Context) {
	Response{code, empty, code.String()}.Json(c)
}

func FailWithCodeAndMsg(code Code, msg string, c *gin.Context) {
	Response{code, empty, msg}.Json(c)
}

func FailWithError(err error, c *gin.Context) {
	data, msg := validata.ValidateError(err)
	FailWithData(data, msg, c)
}

func (c Code) ToResp(data any) Response {
	return Response{Code: c, Data: data, Msg: c.String()}
}
