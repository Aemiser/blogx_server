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
	LogRouter(nr)
	SiteRouter(nr)
	IamgeRouter(nr)
	BannerRouter(nr)
	CaptchaRouter(nr)
	addr := global.Config.System.Addr()
	r.Run(addr)
}
