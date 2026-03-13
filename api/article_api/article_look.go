package article_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/redis_service/redis_article"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type ArticleLookRequest struct {
	ArticleID  uint `json:"articleID" binding:"required"`
	TimeSecond int  `json:"timeSecond"` //读文章一共用了多少时间
}

func (ArticleApi) ArticleLookView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleLookRequest](c)

	// TODO:未登录的用户，浏览量怎么增加
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.SuccessWithData("登入成功", c)
		return
	}

	var article models.ArticleModel
	// 检查文章是否存在，判断文章状态
	err = global.Db.Take(&article, "status = ? and id =?", enum.ArticlePublished, cr.ArticleID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	if redis_article.GetUserArticleHistoryCache(article.ID, claims.Claims.UserID) {
		logrus.Infof("在缓存中")
		res.SuccessWithMsg("成功", c)
		return
	}

	// 查这个文章今天有没有在足迹里面
	var history models.UserArticleLookHistoryModel
	err = global.Db.Take(&history, "article_id = ? and user_id = ? and created_at > ? and created_at < ?",
		cr.ArticleID,
		claims.Claims.UserID,
		time.Now().Format("2006-01-02 15:04:05")+" 00:00:00",
		time.Now().Format("2006-01-02")).Error
	if err == nil {
		res.SuccessWithData("成功", c)
		return
	}

	err = global.Db.Create(&models.UserArticleLookHistoryModel{
		UserID:    claims.Claims.UserID,
		ArticleID: cr.ArticleID,
	}).Error
	if err != nil {
		res.FailWithMsg("失败", c)
		return
	}

	redis_article.SetCacheLook(article.ID, true)
	redis_article.SetUserArticleHistoryCache(article.ID, claims.Claims.UserID)
	res.SuccessWithMsg("成功", c)
	return
}

type ArticleLookListRequest struct {
	common.PageInfo
	UserID uint `form:"userID"`
	Type   int8 `form:"type" binding:"required,oneof=1 2"` // 1为用户 2 为管理员
}

type ArticleLookListResponse struct {
	ID        uint      `json:"id"`
	LookData  time.Time `json:"lookData"`
	Title     string    `json:"title"`
	Cover     string    `json:"cover"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	UserID    uint      `json:"userID"`
	ArticleID uint      `json:"articleID"`
}

func (ArticleApi) ArticleLookListView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleLookListRequest](c)
	claims := jwts.GetClaimsByGin(c)

	switch cr.Type {
	case 1:
		cr.UserID = claims.Claims.UserID
	}

	_list, count, _ := common.ListQuery(models.UserArticleLookHistoryModel{
		UserID: cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Preloads: []string{"UserModel", "ArticleModel"},
	})

	var list = make([]ArticleLookListResponse, 0)
	for _, model := range _list {
		list = append(list, ArticleLookListResponse{
			ID:        model.ID,
			LookData:  model.ArticleModel.CreatedAt,
			Title:     model.ArticleModel.Title,
			Cover:     model.ArticleModel.Cover,
			Nickname:  model.UserModel.Nickname,
			Avatar:    model.UserModel.Avatar,
			UserID:    model.UserID,
			ArticleID: model.ArticleID,
		})
	}
	res.SuccessWithList(list, count, c)
}
