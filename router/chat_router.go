package router

import (
	"blogx_server/api"
	"blogx_server/api/chat_api"
	"blogx_server/middlerware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

func ChatRouter(r *gin.RouterGroup) {
	app := api.App.ChatApi
	r.GET("chat/record", middlerware.AuthMiddleware, middlerware.BindQueryMiddlerware[chat_api.ChatListRequest], app.ChatListView)
	r.GET("chat/session", middlerware.AuthMiddleware, middlerware.BindQueryMiddlerware[chat_api.SessionListRequest], app.SessionListView)
	r.DELETE("chat", middlerware.AuthMiddleware, middlerware.BindJsonMiddlerware[models.IDListRequest], app.UserChatDeleteView)
	r.DELETE("chat/user/:id", middlerware.AuthMiddleware, middlerware.BindUriMiddlerware[models.IDRequest], app.UserChatDeleteSessionChatView)
	r.POST("chat/read/:id", middlerware.AuthMiddleware, middlerware.BindUriMiddlerware[models.IDRequest], app.UserChatReadView)
	r.GET("chat/ws", app.ChatView)
}
