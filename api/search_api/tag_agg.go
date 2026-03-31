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
		return
	}

	agg := elastic.NewTermsAggregation().Field("tag_list")
	agg.SubAggregation("page",
		elastic.NewBucketSortAggregation().
			From(cr.GetOffset()).
			Size(cr.Limit))

	query := elastic.NewBoolQuery()
	//query.MustNot(elastic.NewTermQuery("tag_list", ""))
	result, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(query).
		Aggregation("tags", agg).
		Aggregation("tags_count", elastic.NewCardinalityAggregation().Field("tag_list")).
		Size(0).Do(context.Background())
	if err != nil {
		logrus.Errorf("查询失败: %v", err)
		res.FailWithMsg("查询失败", c)
	}

	var t AggType
	var val = result.Aggregations["tags"]
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
