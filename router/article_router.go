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

	r.POST("article/examine", middlerware.AdminMiddleware, middlerware.BindJsonMiddlerware[article_api.ArticleExamineRequest], app.ArticleExamineView)
	r.GET("article/digg/:id", middlerware.AuthMiddleware, middlerware.BindUriMiddlerware[models.IDRequest], app.ArticleDiggView)
	r.POST("article/collect", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[article_api.ArticleCollectRequest], app.ArticleCollectView)

	r.DELETE("article/:id", middlerware.AuthMiddleware, middlerware.BindUriMiddlerware[models.IDRequest], app.ArticleRemoveUserView)
	r.DELETE("article", middlerware.AdminMiddleware, middlerware.BindJsonMiddlerware[models.IDListRequest], app.ArticleRemoveView)

	r.POST("history", middlerware.BindJsonMiddlerware[article_api.ArticleLookRequest], app.ArticleLookView)
	r.GET("history", middlerware.AuthMiddleware, middlerware.BindQueryMiddlerware[article_api.ArticleLookListRequest], app.ArticleLookListView)
	r.DELETE("history", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[models.IDListRequest], app.ArticleLookDeleteView)

	r.POST("category", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[article_api.CategoryCreate], app.CategoryCreateView)
	r.GET("category", middlerware.BindQueryMiddlerware[article_api.CategoryListRequest], app.CategoryListView)
	r.DELETE("category", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[models.IDListRequest], app.CategoryRemoveView)

	r.POST("collect", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[article_api.CollectCreate], app.CollectCreateView)
	r.GET("collect", middlerware.BindQueryMiddlerware[article_api.CollectListRequest], app.CollectListView)
	r.DELETE("collect", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[models.IDListRequest], app.CollectRemoveView)

	r.GET("category/options", middlerware.AuthMiddleware, app.CategoryOptionsView)

}
