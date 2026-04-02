package router

import (
	"blogx_server/api"
	"blogx_server/api/focus_api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func FocusRouter(r *gin.RouterGroup) {
	app := api.App.FocusApi
	// 关注对方
	r.POST("focus", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[focus_api.FocusUserRequest], app.FocusUserApi)
	// 取关对方
	r.DELETE("focus", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[focus_api.FocusUserRequest], app.UnFocusUserApi)
	// 获取关注列表
	r.GET("focus/my_focus", middlerware.BindQueryMiddlerware[focus_api.FocusUserListRequest], app.FocusUserListApi)
	// 获取粉丝列表
	r.GET("focus/my_fans", middlerware.BindQueryMiddlerware[focus_api.FocusUserListRequest], app.FansUserListApi)
}
