package user_api

import (
	"blogx_server/common/res"
	"blogx_server/service/redis_service/redis_jwt"

	"github.com/gin-gonic/gin"
)

func (UserApi) LogoutView(c *gin.Context) {
	token := c.Request.Header.Get("token")
	redis_jwt.BlackToken(token, redis_jwt.UserBlackType)

	res.SuccessWithMsg("退出成功", c)

}
