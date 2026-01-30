package router

import (
	"blogx_server/api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func LogRouter(r *gin.RouterGroup) {
	/*app := api.App.LogApi
	r.Use(middlerware.AdminMiddleware)
	r.GET("logs", app.LogListView)
	r.GET("logs/:id", app.LogReadView)
	r.DELETE("logs", app.LogRemoveView)*/
	app := api.App.LogApi

	// 为日志相关路由创建单独的子路由组
	logGroup := r.Group("logs")
	logGroup.Use(middlerware.AdminMiddleware)

	logGroup.GET("", app.LogListView)
	logGroup.GET(":id", app.LogReadView)
	logGroup.DELETE("", app.LogRemoveView)
}
