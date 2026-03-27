package router

import (
	"blogx_server/api"
	"blogx_server/api/ai_api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func AiRouter(r *gin.RouterGroup) {
	app := api.App.AiApi
	r.POST("ai/analysis", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[ai_api.AIAnalysisRequest], app.AIAnalysisView)
	r.POST("ai/article", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[ai_api.ArticleAiRequest], app.ArticleAiView)
}
