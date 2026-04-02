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
	"blogx_server/utils/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ArticleListRequest struct {
	common.PageInfo
	Type       int8  `form:"type" binding:"required,oneof=1 2 3 "` // 1看别人的 2看自己的 3管理员看
	UserID     uint  `form:"userID"`
	CategoryID *uint `form:"categoryID"`
	Status     enum.ArticleStatus
	CollectID  uint `form:"collectID"`
}

type ArticleListResponse struct {
	models.ArticleModel
	UserTop       bool    `json:"userTop"`  // 用户是否置顶
	AdminTop      bool    `json:"adminTop"` //管理员是否置顶
	UserNickname  string  `json:"nickName"`
	UserAvatar    string  `json:"userAvatar"`
	CategoryTitle *string `json:"categoryTitle"`
}

func (ArticleApi) ArticleListView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleListRequest](c)

	var topArticleIDList []uint // [1,2,3] => (1,2,3)
	// 排序白名单字段
	var OrderColumnMap = map[string]bool{
		"look_count desc":    true,
		"digg_count desc":    true,
		"comment_count desc": true,
		"collect_count desc": true,
		"look_count asc":     true,
		"digg_count asc":     true,
		"comment_count asc":  true,
		"collect_count asc":  true,
	}

	switch cr.Type {
	case 1:
		// 查别人 用户ID 必填
		if cr.UserID == 0 {
			res.FailWithMsg("用户ID不能为空", c)
			return
		}
		// 查询受限
		if cr.Page > 2 || cr.Limit > 10 {
			res.FailWithMsg("查询更多，请登入", c)
			return
		}
		cr.Status = 0
		cr.Order = ""

		if cr.CollectID != 0 {
			// 如果传入了收藏夹ID，查权限
			var userconf models.UserConfigModel
			err := global.Db.Take(&userconf, "user_id = ?", cr.UserID).Error
			if err != nil {
				res.FailWithMsg("用户不存在", c)
				return
			}

			if !userconf.OpenCollect {
				res.FailWithMsg("用户未开放收藏功能", c)
				return
			}
		}
	case 2:
		// 查自己
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			res.FailWithMsg("请登录", c)
			return
		}
		cr.UserID = claims.Claims.UserID
	case 3:
		// 管理员
		claims, err := jwts.ParseTokenByGin(c)
		if !(err == nil && claims.Claims.Role == enum.AdminRole) {
			res.FailWithMsg("角色错误", c)
			return
		}
	}
	query := global.Db.Where("")
	if cr.CollectID != 0 {
		fmt.Println("收藏夹ID", cr.CollectID)
		var articleIDList []uint
		global.Db.Model(models.UserArticleCollectModel{}).Where("collect_id = ?", cr.CollectID).Select("article_id").Scan(&articleIDList)
		fmt.Println("收藏夹IDList", articleIDList)
		query.Where("id in ?", articleIDList)
	}

	// 对于类型2,3而言存在order判断
	if cr.Order != "" {
		_, ok := OrderColumnMap[cr.Order]
		if !ok {
			res.FailWithMsg("排序字段错误", c)
			return
		}
	}

	// 置顶文章处理
	var userTopMap = map[uint]bool{}
	var AdminTopMap = map[uint]bool{}
	if cr.UserID != 0 {
		var userTopArticleList []models.UserTopArticleModel
		// 这里查询出来的置顶文章按照时间升序，在下面的文章列表查询中，即可把这些置顶的文章按照升序排列到最前面
		global.Db.Preload("UserModel").Order("created_at").Take(&userTopArticleList, "user_id = ?", cr.UserID)

		for _, i2 := range userTopArticleList {
			topArticleIDList = append(topArticleIDList, i2.ArticleID)
			if i2.UserModel.Role == enum.AdminRole {
				AdminTopMap[i2.ArticleID] = true
			}
			userTopMap[i2.ArticleID] = true
		}
	}

	// 判断topArticleIDList为空的情况
	var option = common.Options{
		Likes:        []string{"title"},
		PageInfo:     cr.PageInfo,
		DefaultOrder: "created_at desc",
		Where:        query,
		Preloads:     []string{"Category", "UserModel"},
		Debug:        true,
	}

	if len(topArticleIDList) > 0 {
		// 这里通过函数，把[3,4,5] => string类型的  id = 3 desc,id = 4 desc,id = 5 desc,这里置顶的顺序按照最新时间降序，在上面获取topArticleIDList列表的sql语句中已经实现了
		option.DefaultOrder = fmt.Sprintf("%s,created_at desc", sql.CoverSliceOrderSql(topArticleIDList))
	}

	// 文章列表查询
	_list, count, _ := common.ListQuery(models.ArticleModel{
		UserID:     cr.UserID,
		CategoryID: cr.CategoryID,
		Status:     cr.Status,
	}, option)
	// 响应数据封装
	var list = make([]ArticleListResponse, 0)
	collectMap := redis_article.GetAllCacheCollect()
	lookMap := redis_article.GetAllCacheLook()
	diggMap := redis_article.GetAllCacheDigg()
	commentMap := redis_article.GetAllCacheComment()

	for _, model := range _list {
		model.Content = ""
		model.DiggCount = model.DiggCount + diggMap[model.ID]
		model.CollectCount = model.CollectCount + collectMap[model.ID]
		model.LookCount = model.LookCount + lookMap[model.ID]
		model.CommentCount = model.LookCount + commentMap[model.ID]
		date := ArticleListResponse{
			ArticleModel: model,
			UserTop:      userTopMap[model.ID],
			AdminTop:     AdminTopMap[model.ID],
			UserNickname: model.UserModel.Nickname,
			UserAvatar:   model.UserModel.Avatar,
		}

		if model.Category != nil {
			date.CategoryTitle = &model.Category.Title
		}
		list = append(list, date)
	}
	res.SuccessWithList(list, count, c)
}
