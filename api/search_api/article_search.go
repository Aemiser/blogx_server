package search_api

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
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type ArticleSearchRequest struct {
	common.PageInfo
	Tag  string `form:"tag"`
	Type int8   `form:"type"` // 0 猜你喜欢 1 最新发布 2最多回复 3最多点赞 4最多收藏
}

type ArticleBaseInfo struct {
	Id       uint   `json:"id"`
	Title    string `json:"title"`
	Abstract string `json:"abstract"`
}

type ArticleListResponse struct {
	models.ArticleModel
	AdminTop      bool    `json:"adminTop"` //管理员是否置顶
	UserNickname  string  `json:"nickName"`
	UserAvatar    string  `json:"userAvatar"`
	CategoryTitle *string `json:"categoryTitle"`
}

func (SearchApi) ArticleSearchView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleSearchRequest](c)
	var searchArticleMap = map[uint]ArticleBaseInfo{}
	var articleIDList []uint
	sortMap := map[int8]string{
		0: "_score",
		1: "created_at",
		2: "comment_count",
		3: "digg_count",
		4: "collect_count",
	}
	sortKey := sortMap[cr.Type]
	if sortKey == "" {
		res.FailWithMsg("搜索类型错误", c)
		return
	}

	// 缓存
	collectMap := redis_article.GetAllCacheCollect()
	lookMap := redis_article.GetAllCacheLook()
	diggMap := redis_article.GetAllCacheDigg()
	commentMap := redis_article.GetAllCacheComment()

	topArticleIDList := getAdminTopArticleIDList()
	// 服务降级
	if global.ESClient == nil {
		where := global.Db.Where("")
		if cr.Tag != "" {
			where.Where("tag_list like ?", fmt.Sprintf("%%%s%%", cr.Tag))
		}

		var articleTopMap = map[uint]bool{}
		for _, u := range topArticleIDList {
			articleTopMap[u] = true
		}

		sortMap = map[int8]string{
			0: "",
			1: "created_at desc",
			2: "comment_count desc",
			3: "digg_count desc",
			4: "collect_count desc",
		}
		sort, _ := sortMap[cr.Type]
		cr.Order = sort
		_list, count, _ := common.ListQuery(models.ArticleModel{}, common.Options{
			Preloads:     []string{"Category", "UserModel"},
			PageInfo:     cr.PageInfo,
			Likes:        []string{"title", "abstract"},
			DefaultOrder: sql.CoverSliceOrderSql(topArticleIDList),
			Where:        where,
		})

		var list = make([]ArticleListResponse, 0)
		for _, model := range _list {
			model.Content = ""
			model.DiggCount = model.DiggCount + diggMap[model.ID]
			model.CollectCount = model.CollectCount + collectMap[model.ID]
			model.LookCount = model.LookCount + lookMap[model.ID]
			model.CommentCount = model.LookCount + commentMap[model.ID]
			item := ArticleListResponse{
				ArticleModel: model,
				AdminTop:     articleTopMap[model.ID],
				UserNickname: model.UserModel.Nickname,
				UserAvatar:   model.UserModel.Avatar,
			}
			if model.Category != nil {
				item.CategoryTitle = &model.Category.Title
			}
			list = append(list, item)
		}
		res.SuccessWithList(list, count, c)
		return
	}

	query := elastic.NewBoolQuery()
	if cr.Key != "" {
		query.Should(
			elastic.NewMatchQuery("title", cr.Key),
			elastic.NewMatchQuery("abstract", cr.Key),
			elastic.NewMatchQuery("content", cr.Key),
		)
	}

	// 把管理员置顶的文章查出来

	var articleTopMap = map[uint]bool{}
	if len(topArticleIDList) > 0 {
		var topArticleIDListAny []interface{}
		for _, u := range topArticleIDList {
			topArticleIDListAny = append(topArticleIDListAny, u)
			articleTopMap[u] = true
			// 不管什么分类都在最上面
			articleIDList = append(articleIDList, u)
		}
		query.Should(elastic.NewTermsQuery("id", topArticleIDListAny...).Boost(10))
	}

	if cr.Type == 0 { // 推荐
		claims, err := jwts.ParseTokenByGin(c)
		if err == nil && claims != nil {
			// 用户登录了
			// 查用户感兴趣的分类
			var userConf models.UserConfigModel
			err = global.Db.Take(&userConf, "user_id = ?", claims.Claims.UserID).Error
			if err != nil {
				res.FailWithMsg("用户配置不存在", c)
				return
			}

			if len(userConf.LikeTags) > 0 {
				tagQuery := elastic.NewBoolQuery()
				var tagAnyList []interface{}
				for _, tag := range userConf.LikeTags {
					tagAnyList = append(tagAnyList, tag)
				}
				tagQuery.Should(elastic.NewTermsQuery("tag_list", tagAnyList...))
				query.Must(tagQuery)
			}
		}
	}
	if cr.Tag != "" {
		query.Must(elastic.NewTermQuery("tag_list", cr.Tag))
	}
	// 只查能发布的文章
	query.Must(elastic.NewTermQuery("status", 3))
	highligth := elastic.NewHighlight()
	highligth.Field("title")
	result, err := global.ESClient.Search(models.ArticleModel{}.Index()).
		Query(query).
		From(cr.GetOffset()).
		Size(cr.GetLimit()).
		Highlight(highligth).
		Sort(sortKey, false).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	count := result.Hits.TotalHits.Value // 总数

	fmt.Printf("总数：%d, 当前页偏移：%d, 每页大小：%d, ES 返回的 Hits 数量：%d\n",
		count, cr.GetOffset(), cr.GetLimit(), len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		fmt.Println(string(hit.Source))
		var art ArticleBaseInfo
		err = json.Unmarshal(hit.Source, &art)
		if err != nil {
			logrus.Warnf("解析失败: %s %s ", err, string(hit.Source))
			continue
		}
		fmt.Println(hit.Highlight["title"])
		if len(hit.Highlight["title"]) > 0 {
			art.Title = hit.Highlight["title"][0]
		}

		if len(hit.Highlight["abstract"]) > 0 {
			art.Abstract = hit.Highlight["abstract"][0]
		}

		searchArticleMap[art.Id] = art
		articleIDList = append(articleIDList, art.Id)
	}
	fmt.Println(articleIDList)
	where := global.Db.Where("id in ?", articleIDList)
	_list, _, _ := common.ListQuery(models.ArticleModel{}, common.Options{
		Where:        where,
		Preloads:     []string{"Category", "UserModel"},
		DefaultOrder: sql.CoverSliceOrderSql(articleIDList),
	})

	var list = make([]ArticleListResponse, 0)
	for _, model := range _list {
		model.Content = ""
		model.DiggCount = model.DiggCount + diggMap[model.ID]
		model.CollectCount = model.CollectCount + collectMap[model.ID]
		model.LookCount = model.LookCount + lookMap[model.ID]
		model.CommentCount = model.LookCount + commentMap[model.ID]
		item := ArticleListResponse{
			ArticleModel: model,
			AdminTop:     articleTopMap[model.ID],
			UserNickname: model.UserModel.Nickname,
			UserAvatar:   model.UserModel.Avatar,
		}
		if model.Category != nil {
			item.CategoryTitle = &model.Category.Title
		}
		item.Title = searchArticleMap[model.ID].Title
		item.Abstract = searchArticleMap[model.ID].Abstract
		list = append(list, item)
	}
	res.SuccessWithList(list, int(count), c)
}

func getAdminTopArticleIDList() (topArticleIDList []uint) {
	var userIDList []uint
	global.Db.Model(models.UserModel{}).Where("role = ?", enum.AdminRole).Select("id").Scan(&userIDList)
	global.Db.Model(models.UserTopArticleModel{}).Where("user_id in ?", userIDList).Select("article_id").Scan(&topArticleIDList)
	return
}
