package middlerware

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/models/enum"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Set("claims", claims)
}

func AdminMiddleware(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	if claims.Claims.Role != enum.AdminRole {
		res.FailWithMsg("权限不足", c)
		c.Abort()
		return
	}
	c.Set("claims", claims)
}
