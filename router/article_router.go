package router

import (
	"blogx_server/api"
	"blogx_server/api/article_api"
	"blogx_server/common"
	"blogx_server/middlerware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

func ArticleRouter(r *gin.RouterGroup) {
	app := api.App.ArticleApi
	// 发布文章
	r.POST("article", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[article_api.ArticleCreateRequest], app.ArticleCreateView)
	r.PUT("article", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[article_api.ArticleUpdateRequest], app.ArticleUpdateView)
	r.GET("article", middlerware.BindQueryMiddlerware[article_api.ArticleListRequest], app.ArticleListView)
	r.GET("article/:id", middlerware.BindUriMiddlerware[models.IDRequest], app.ArticleDetailView)

	// 审核点赞
	r.POST("article/examine", middlerware.AdminMiddleware, middlerware.BindJsonMiddlerware[article_api.ArticleExamineRequest], app.ArticleExamineView)
	r.GET("article/digg/:id", middlerware.AuthMiddleware, middlerware.BindUriMiddlerware[models.IDRequest], app.ArticleDiggView)

	// 删除文章
	r.DELETE("article/:id", middlerware.AuthMiddleware, middlerware.BindUriMiddlerware[models.IDRequest], app.ArticleRemoveUserView)
	r.DELETE("article", middlerware.AdminMiddleware, middlerware.BindJsonMiddlerware[models.IDListRequest], app.ArticleRemoveView)

	// 浏览记录
	r.POST("history", middlerware.BindJsonMiddlerware[article_api.ArticleLookRequest], app.ArticleLookView)
	r.GET("history", middlerware.AuthMiddleware, middlerware.BindQueryMiddlerware[article_api.ArticleLookListRequest], app.ArticleLookListView)
	r.DELETE("history", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[models.IDListRequest], app.ArticleLookDeleteView)

	// 分类
	r.POST("category", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[article_api.CategoryCreate], app.CategoryCreateView)
	r.GET("category", middlerware.BindQueryMiddlerware[article_api.CategoryListRequest], app.CategoryListView)
	r.DELETE("category", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[models.IDListRequest], app.CategoryRemoveView)

	// 收藏文章
	r.POST("article/collect", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[article_api.ArticleCollectRequest], app.ArticleCollectView)
	r.DELETE("article/collect", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[article_api.ArticleCollectRemoveRequest], app.ArticleCollectRemoveView)

	// 收藏夹
	r.POST("collect", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[article_api.CollectCreate], app.CollectCreateView)
	r.GET("collect", middlerware.BindQueryMiddlerware[article_api.CollectListRequest], app.CollectListView)
	r.DELETE("collect", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[models.IDListRequest], app.CollectRemoveView)

	//标签
	r.GET("category/options", middlerware.AuthMiddleware, app.CategoryOptionsView)
	r.GET("article/tag/options", middlerware.AuthMiddleware, app.ArticleTagListView)

	// 推荐
	r.GET("article/auth_recommend", middlerware.BindQueryMiddlerware[common.PageInfo], app.AuthRecommentView)
	r.GET("article/article_recommend", middlerware.BindQueryMiddlerware[common.PageInfo], app.ArticleRecommentView)

	//添加到足迹
	r.POST("article/look", middlerware.BindJsonMiddlerware[article_api.ArticleLookRequest], app.ArticleLookView)
}
