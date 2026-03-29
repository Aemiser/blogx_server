package router

import (
	"blogx_server/api"
	"blogx_server/api/user_api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	app := api.App.UserApi
	r.POST("user/send_email", middlerware.CaptchaMiddleware, app.SendEmailView)
	r.POST("user/email", app.RegisterEmail)
	r.POST("user/qq", app.QQLoginView)
	r.POST("user/login", middlerware.CaptchaMiddleware, middlerware.BindJsonMiddlerware[user_api.PwdLoginRequest], app.PwdLoginApi)
	r.GET("user/detail", middlerware.AuthMiddleware, app.UserDetailView)
	r.GET("user", middlerware.AdminMiddleware, middlerware.BindQueryMiddlerware[user_api.UserListRequest], app.UserListView)
	r.GET("user/login", middlerware.AuthMiddleware, app.UserLoginListView)
	r.GET("user/base", app.UserBaseInfoView)
	r.PUT("user/password", middlerware.AuthMiddleware, app.UpdatePasswordView)
	r.PUT("user/password/reset", middlerware.EmailVerifyMiddleware, app.ResetPassowrdView)
	r.PUT("user/email/bind", middlerware.EmailVerifyMiddleware, middlerware.AuthMiddleware, app.BindEmail)
	r.PUT("user", middlerware.AuthMiddleware, app.UserInfoUpdate)
	r.PUT("user/admin", middlerware.AdminMiddleware, app.AdminUserInfoUpdate)
}
