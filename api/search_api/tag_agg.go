package search_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type AggType struct {
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
	Buckets                 []struct {
		Key      string `json:"key"`
		DocCount int    `json:"doc_count"`
	} `json:"buckets"`
}

type TagAggResponse struct {
	Tag          string `json:"tag"`
	ArticleCount int    `json:"articleCount"`
}

type AggCount struct {
	Value int `json:"value"`
}

func (SearchApi) TagAggView(c *gin.Context) {
	var cr = middlerware.GetBind[common.PageInfo](c)

	var list = make([]TagAggResponse, 0)
	if global.ESClient == nil {
		var articleList []models.ArticleModel
		global.Db.Find(&articleList, "tag_list <> ''")

		var tagMap = map[string]int{}
		for _, model := range articleList {
			for _, tag := range model.TagList {
				count, ok := tagMap[tag]
				if !ok {
					tagMap[tag] = 1
					continue
				}
				tagMap[tag] = count + 1
			}
		}
		for tag, count := range tagMap {
			list = append(list, TagAggResponse{
				Tag:          tag,
				ArticleCount: count,
			})
		}
		if cr.Limit >= len(list) {
			cr.Limit = len(list)
		}
		res.SuccessWithList(list[0:cr.Limit], len(list), c)
		return
	}

	// 不使用 bucket_sort 分页，直接在 terms 聚合中设置 size
	agg := elastic.NewTermsAggregation().Field("tag_list").Size(cr.GetLimit())

	query := elastic.NewBoolQuery()
	//query.MustNot(elastic.NewTermQuery("tag_list", ""))
	result, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(query).
		Aggregation("tags", agg).
		Aggregation("tags_count", elastic.NewCardinalityAggregation().Field("tag_list")).
		Size(0).Do(context.Background())
	if err != nil {
		logrus.Errorf("查询失败：%v", err)
		res.FailWithMsg("查询失败", c)
		return
	}

	// 检查 result 是否为 nil
	if result == nil || result.Aggregations == nil {
		logrus.Errorf("ES 返回结果为空")
		res.FailWithMsg("查询失败", c)
		return
	}

	var t AggType
	var val = result.Aggregations["tags"]
	if val == nil {
		logrus.Errorf("tags 聚合结果为空")
		res.FailWithMsg("查询失败", c)
		return
	}
	err = json.Unmarshal(val, &t)
	if err != nil {
		logrus.Errorf("解析json失败:%s %s", err, string(val))
		res.FailWithMsg("查询失败", c)
	}

	var co AggCount
	err = json.Unmarshal(result.Aggregations["tags_count"], &co)
	if err != nil {
		logrus.Errorf("解析json失败:%s %s", err, string(result.Aggregations["tags_count"]))
		res.FailWithMsg("查询失败", c)
	}

	for _, bucket := range t.Buckets {
		list = append(list, TagAggResponse{
			Tag:          bucket.Key,
			ArticleCount: bucket.DocCount,
		})
	}
	res.SuccessWithList(list, co.Value, c)

}
