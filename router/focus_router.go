package router

import (
	"blogx_server/api"
	"blogx_server/api/focus_api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func FocusRouter(r *gin.RouterGroup) {
	app := api.App.FocusApi
	r.POST("focus", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[focus_api.FocusUserRequest], app.FocusUserApi)
	r.DELETE("focus", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[focus_api.FocusUserRequest], app.UnFocusUserApi)
	r.GET("focus/my_focus", middlerware.BindQueryMiddlerware[focus_api.FocusUserListRequest], app.FocusUserListApi)
	r.GET("focus/my_fans", middlerware.BindQueryMiddlerware[focus_api.FocusUserListRequest], app.FansUserListApi)
}
