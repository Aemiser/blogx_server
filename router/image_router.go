package router

import (
	"blogx_server/api"
	"blogx_server/middlerware"

	"github.com/gin-gonic/gin"
)

func IamgeRouter(r *gin.RouterGroup) {
	app := api.App.Image
	r.POST("images", middlerware.AuthMiddleware, app.ImageUploadView)
	r.GET("images", middlerware.AdminMiddleware, app.ImageListView)
	r.DELETE("images", middlerware.AdminMiddleware, app.ImageRemoveView)

}
