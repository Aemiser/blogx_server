package middlerware

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LogMiddleware(c *gin.Context) {
	// 请求中间件
	byteData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Error(err.Error())
	}
	fmt.Println("boby: ", string(byteData))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(byteData))

	c.Next()
	// 响应中间件
}
