package router

import (
	"blogx_server/api"
	"blogx_server/api/article_api"
	"blogx_server/middlerware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

func ArticleRouter(r *gin.RouterGroup) {
	app := api.App.ArticleApi
	r.POST("article", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[article_api.ArticleCreateRequest], app.ArticleCreateView)
	r.PUT("article", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[article_api.ArticleUpdateRequest], app.ArticleUpdateView)
	r.GET("article", middlerware.BindQueryMiddlerware[article_api.ArticleListRequest], app.ArticleListView)
	r.GET("article/:id", middlerware.BindUriMiddlerware[models.IDRequest], app.ArticleDetailView)
	r.POST("article/examine", middlerware.BindJsonMiddlerware[article_api.ArticleExamineRequest], app.ArticleExamineView)
}
