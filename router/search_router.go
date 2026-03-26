package router

import (
	"blogx_server/api"
	"blogx_server/api/search_api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func SearchRouter(r *gin.RouterGroup) {
	app := api.App.SearchApi
	r.GET("search/article", middlerware.BindQueryMiddlerware[search_api.ArticleSearchRequest], app.ArticleSearchView)
}
