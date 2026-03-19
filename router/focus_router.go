package router

import (
	"blogx_server/api"
	"blogx_server/api/focus_api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func FocusRouter(r *gin.RouterGroup) {
	app := api.App.FocusApi
	r.POST("focus", middlerware.AdminMiddleware, middlerware.BindJsonMiddlerware[focus_api.FocusUserRequest], app.FocusUserApi)
	r.GET("focus/my_focus", middlerware.AuthMiddleware, middlerware.BindQueryMiddlerware[focus_api.FocusUserListRequest], app.FocusUserListApi)
}
