package router

import (
	"blogx_server/api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	app := api.App.UserApi
	r.POST("user/send_email", middlerware.CaptchaMiddleware, app.SendEmailView)

}
