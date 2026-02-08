package user_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"

	"github.com/gin-gonic/gin"
)

func (UserApi) BindEmail(c *gin.Context) {
	if !global.Config.Site.Login.EmailLogin {
		res.FailWithMsg("站点未启用邮箱注册", c)
		return
	}

	_email, _ := c.Get("email")
	email := _email.(string)

	user, err := jwts.GetClaimsByGin(c).GetUser()
	if err != nil {
		res.FailWithMsg("不存在用户", c)
		return
	}

	err = global.Db.Model(&user).Update("email", email).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.SuccessWithMsg("绑定成功", c)

}
