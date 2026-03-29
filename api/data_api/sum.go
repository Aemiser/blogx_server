package data_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/redis_service/redis_site"

	"github.com/gin-gonic/gin"
)

type SumResponse struct {
	FlowCount     int   `json:"flowCount"`
	UserCount     int64 `json:"userCount"`
	ArticleCount  int64 `json:"articleCount"`
	MessageCount  int64 `json:"messageCount"`
	CommentCount  int64 `json:"commentCount"`
	NewLoginCount int64 `json:"newLoginCount"`
	NewSignCount  int64 `json:"newSignCount"`
}

func (DataApi) SumView(c *gin.Context) {
	var data SumResponse
	data.FlowCount = redis_site.GetFlow()
	global.Db.Model(models.UserModel{}).Count(&data.UserCount)
	global.Db.Model(models.ArticleModel{}).Where("status = ?", enum.ArticlePublished).Count(&data.ArticleCount)
	global.Db.Model(models.ChatModel{}).Count(&data.MessageCount)
	global.Db.Model(models.CommentModel{}).Count(&data.CommentCount)
	global.Db.Model(models.UserLoginModel{}).Where("date(created_at) = date(now())").Count(&data.NewLoginCount)
	global.Db.Model(models.UserModel{}).Where("date(created_at) = date(now())").Count(&data.NewSignCount)
	res.SuccessWithData(data, c)
}
