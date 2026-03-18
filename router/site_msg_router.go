package router

import (
	"blogx_server/api"
	"blogx_server/api/site_msg_api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func SiteMsgRouter(r *gin.RouterGroup) {
	app := api.App.SiteMsgAPi
	r.GET("site_msg/conf", middlerware.AuthMiddleware, app.UserSiteMessageConfView)
	r.PUT("site_msg/conf", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[site_msg_api.UserMessageConfModelRequest], app.UserSiteMessageUpdateConfView)
}
