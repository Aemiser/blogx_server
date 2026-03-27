package res

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func SSEOK(data any, c *gin.Context) {
	byteData, _ := json.Marshal(Response{SuccessCode, data, "成功"})
	c.SSEvent("", string(byteData))
	c.Writer.Flush()
}

func SSEFail(data any, c *gin.Context) {
	byteData, _ := json.Marshal(Response{FailServiceCode, map[string]any{}, "成功"})
	c.SSEvent("", string(byteData))
	c.Writer.Flush()
}
