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
	addr := global.Config.System.Addr()
	r.Run(addr)
}
