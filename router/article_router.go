package router

import (
	"blogx_server/api"
	"blogx_server/api/article_api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func ArticleRouter(r *gin.RouterGroup) {
	app := api.App.ArticleApi
	r.POST("article", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[article_api.ArticleCreateRequest], app.ArticleCreateView)
}
