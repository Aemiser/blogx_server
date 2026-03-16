package router

import (
	"blogx_server/api"
	"blogx_server/api/comment_api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func CommentRouter(r *gin.RouterGroup) {
	app := api.App.CommentApi
	r.POST("comment", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[comment_api.CommentCreateReaquest], app.CommentListView)
}
