package article_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/redis_service/redis_article"

	"github.com/gin-gonic/gin"
)

type ArticleDetailResponse struct {
	models.ArticleModel
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	UserAvatar string `json:"useravatar"`
}

func (ArticleApi) ArticleDetailView(c *gin.Context) {
	cr := middlerware.GetBind[models.IDRequest](c)

	var article models.ArticleModel
	err := global.Db.Take(&article, cr.ID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}
	//未登入的用户只能看见发布成功的文章

	//登入的用户可以看见自己的全部文章

	//管理员看所有文章
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		// 未登入的用户，查看不是已发步的文章，就返回文章不存在，即未登入的用户只能看见发布成功的文章
		if article.Status != enum.ArticlePublished {
			res.FailWithMsg("文章不存在", c)
			return
		}
		// token 无效时，设置为空 claims，避免后续访问空指针
		claims = nil
	}

	// 如果 claims 为空，说明 token 无效，直接跳过权限判断
	if claims != nil {
		switch claims.Claims.Role {
		case enum.UserRole:
			// 用户只能看见自己的全部文章，但不能看别人未发布的文章
			if article.UserID != claims.Claims.UserID {
				// 不是自己的文章
				if article.Status != enum.ArticlePublished {
					// 文章不是已发布
					res.FailWithMsg("文章不存在", c)
					return
				}
			}
		}
	}

	collentCount := redis_article.GetArticleCacheCollect(article.ID)
	DiggCount := redis_article.GetArticleCacheDigg(article.ID)
	LookCount := redis_article.GetArticleCacheLook(article.ID)

	article.LookCount = article.LookCount + LookCount
	article.DiggCount = article.DiggCount + DiggCount
	article.CollectCount = article.CollectCount + collentCount
	var resp = ArticleDetailResponse{
		ArticleModel: article,
		Username:     article.UserModel.Username,
		Nickname:     article.UserModel.Nickname,
		UserAvatar:   article.UserModel.Avatar,
	}

	res.SuccessWithData(resp, c)
}
