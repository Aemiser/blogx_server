package router

import (
	"blogx_server/api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func BannerRouter(r *gin.RouterGroup) {
	app := api.App.Banner
	r.POST("banner", middlerware.AdminMiddleware, app.BannerCreateView)
	r.DELETE("banner", middlerware.AdminMiddleware, app.BannerRemoveView)
	r.PUT("banner/:id", middlerware.AdminMiddleware, app.BannerUpdateView)
	r.GET("banner", middlerware.AuthMiddleware, app.BannerListView)
}
