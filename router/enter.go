package router

import (
	"blogx_server/global"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	nr := r.Group("/api")
	nr.Use(middlerware.LogMiddleware)
	r.Static("/uploads", "uploads")

	SiteRouter(nr)
	LogRouter(nr)
	IamgeRouter(nr)
	BannerRouter(nr)
	CaptchaRouter(nr)
	UserRouter(nr)
	ArticleRouter(nr)
	CommentRouter(nr)
	SiteMsgRouter(nr)
	GlobalNotificationRouter(nr)
	FocusRouter(nr)
	ChatRouter(nr)
	SearchRouter(nr)
	addr := global.Config.System.Addr()
	r.Run(addr)
}
