package router

import (
	"blogx_server/api"
	"blogx_server/api/ai_api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func AiRouter(r *gin.RouterGroup) {
	app := api.App.AiApi
	// 文章分析: 分析文章的标题、摘要、分类、标签
	r.POST("ai/analysis", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[ai_api.AIAnalysisRequest], app.AIAnalysisView)
	// ai智能推荐文章
	r.GET("ai/article", middlerware.AuthMiddleware, middlerware.BindQueryMiddlerware[ai_api.ArticleAiRequest], app.ArticleAiView)
}
