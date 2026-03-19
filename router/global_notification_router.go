package router

import (
	"blogx_server/api"
	"blogx_server/api/global_notification_api"
	"blogx_server/middlerware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

func GlobalNotificationRouter(r *gin.RouterGroup) {
	app := api.App.GlobalNotificationApi
	r.POST("global_notification", middlerware.AdminMiddleware, middlerware.BindJsonMiddlerware[global_notification_api.CreateRequest], app.CreateView)
	r.GET("global_notification", middlerware.AuthMiddleware, middlerware.BindQueryMiddlerware[global_notification_api.ListRequest], app.ListView)
	r.DELETE("global_notification", middlerware.AdminMiddleware, middlerware.BindJsonMiddlerware[models.IDListRequest], app.RemoveView)
	r.POST("global_notification/user", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[global_notification_api.UserMsgActionRequest], app.UserMsgActionView)
}
