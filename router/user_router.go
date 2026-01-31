package router

import (
	"blogx_server/api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	app := api.App.UserApi
	r.POST("user/send_email", middlerware.CaptchaMiddleware, app.SendEmailView)
	r.POST("user/email", app.RegisterEmail)
	r.POST("user/qq", app.QQLoginView)
	r.POST("user/login", middlerware.CaptchaMiddleware, app.PwdLoginApi)
	r.GET("user/detail", middlerware.AuthMiddleware, app.UserDetailView)
	r.GET("user/login", middlerware.AuthMiddleware, app.UserLoginListView)
	r.GET("user/base", app.UserBaseInfoView)
}
