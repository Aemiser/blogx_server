package router

import (
	"blogx_server/api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func DataRouter(r *gin.RouterGroup) {
	app := api.App.DataApi
	r.GET("data/sum", middlerware.AdminMiddleware, app.SumView)
	r.GET("data/article", middlerware.AdminMiddleware, app.ArticleDataView)
}
