package middlerware

import (
	"blogx_server/common/res"

	"github.com/gin-gonic/gin"
)

func BindJsonMiddlerware[T any](c *gin.Context) {
	var req T
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Set("request", req)
}

func BindQueryMiddlerware[T any](c *gin.Context) {
	var req T
	err := c.ShouldBindQuery(&req)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Set("request", req)

}

func GetBind[T any](c *gin.Context) (cr T) {
	return c.MustGet("request").(T)
}
