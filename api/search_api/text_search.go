package search_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/service/text_service"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type TextSearchRequest struct {
	common.PageInfo
}

type TextListResponse struct {
	ArticleID uint   `json:"articleID"`
	Head      string `json:"head"`
	Body      string `json:"body"`
}

func (SearchApi) TextSearchView(c *gin.Context) {
	cr := middlerware.GetBind[TextSearchRequest](c)

	// 服务降级
	if global.ESClient == nil {
		_list, count, _ := common.ListQuery(models.TextModel{}, common.Options{
			PageInfo: cr.PageInfo,
			Likes:    []string{"head", "body"},
		})

		var list = make([]TextListResponse, 0)
		for _, model := range _list {
			list = append(list, TextListResponse{
				ArticleID: model.ArticleID,
				Head:      model.Head,
				Body:      model.Body,
			})
		}
		res.SuccessWithList(list, count, c)
		return
	}

	query := elastic.NewBoolQuery()
	if cr.Key != "" {
		query.Should(
			elastic.NewMatchQuery("head", cr.Key),
			elastic.NewMatchQuery("body", cr.Key),
		)
	}

	// 只查能发布的文章
	highligth := elastic.NewHighlight()
	highligth.Field("head")
	highligth.Field("body")

	result, err := global.ESClient.Search(models.TextModel{}.Index()).
		Query(query).
		From(cr.GetOffset()).
		Size(cr.GetLimit()).
		Highlight(highligth).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	count := result.Hits.TotalHits.Value // 总数

	fmt.Printf("总数：%d, 当前页偏移：%d, 每页大小：%d, ES 返回的 Hits 数量：%d\n",
		count, cr.GetOffset(), cr.GetLimit(), len(result.Hits.Hits))

	var list = make([]TextListResponse, 0)
	for _, hit := range result.Hits.Hits {

		var item text_service.TextModel
		err = json.Unmarshal(hit.Source, &item)
		if err != nil {
			logrus.Warnf("解析失败: %s %s ", err, string(hit.Source))
			continue
		}
		if len(hit.Highlight["head"]) > 0 {
			item.Head = hit.Highlight["head"][0]
		}

		if len(hit.Highlight["body"]) > 0 {
			item.Body = hit.Highlight["body"][0]
		}

		list = append(list, TextListResponse{
			ArticleID: item.ArticleID,
			Head:      item.Head,
			Body:      item.Body,
		})
	}
	res.SuccessWithList(list, int(count), c)
}
