package search_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

type ArticleSearchRequest struct {
	common.PageInfo
	Type int8 `json:"type"` // 0 猜你喜欢 1 最新发布 2最多回复 3最多点赞 4最多收藏
}

func (SearchApi) ArticleSearchView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleSearchRequest](c)

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

	query := elastic.NewBoolQuery()
	if cr.Key != "" {
		query.Should(
			elastic.NewMatchQuery("title", cr.Key),
			elastic.NewMatchQuery("abstract", cr.Key),
			elastic.NewMatchQuery("content", cr.Key),
		)
	}
	claims, err := jwts.ParseTokenByGin(c)
	if err == nil && claims != nil {
		// 用户登录了
		// 查用户感兴趣的分类
		var userConf models.UserConfigModel
		err = global.Db.Take(&userConf, "user_id = ?", claims.Claims.UserID).Error
		if err != nil {
			res.FailWithMsg("用户配置不存在", c)
		}

		if len(userConf.LikeTags) > 0 {
			tagQuery := elastic.NewBoolQuery()
			for _, tag := range userConf.LikeTags {
				tagQuery.Should(elastic.NewTermQuery("tag_list", tag))
			}
			query.Must(tagQuery)
		}
	}
	// 只查能发布的文章
	query.Must(elastic.NewTermQuery("status", 3))
	// 把管理员置顶的文章查出来
	highligth := elastic.NewHighlight()
	highligth.Field("title")
	res, err := global.ESClient.Search(models.ArticleModel{}.Index()).
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
	count := res.Hits.TotalHits.Value // 总数
	fmt.Println(count)
	for _, hit := range res.Hits.Hits {
		fmt.Println(string(hit.Source))
		fmt.Println(hit.Highlight["title"])
	}

}
