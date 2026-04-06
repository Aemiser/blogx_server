package ai_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/ai_service"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

func (AiApi) EmbeddingGenerateView(c *gin.Context) {
	if global.ESClient == nil {
		res.FailWithMsg("ES 未启用", c)
		return
	}

	var articles []models.ArticleModel
	err := global.Db.Where("status = ?", 3).Find(&articles).Error
	if err != nil {
		res.FailWithMsg("查询文章失败", c)
		return
	}

	fmt.Printf("共找到 %d 篇已发布文章\n", len(articles))

	success := 0
	failed := 0

	for i, article := range articles {
		text := article.Title + " " + article.Abstract
		if text == " " {
			continue
		}

		vector, err := ai_service.GetEmbedding(text)
		if err != nil {
			logrus.Errorf("文章 %d embedding 生成失败: %v", article.ID, err)
			failed++
			continue
		}

		_, err = global.ESClient.Update().
			Index(models.ArticleModel{}.Index()).
			Id(fmt.Sprintf("%d", article.ID)).
			Doc(map[string]interface{}{
				"content_vector": vector,
			}).
			Do(context.Background())
		if err != nil {
			logrus.Errorf("文章 %d ES 更新失败: %v", article.ID, err)
			failed++
			continue
		}

		success++
		fmt.Printf("[%d/%d] 文章 %d 完成\n", i+1, len(articles), article.ID)
	}

	res.Success(map[string]interface{}{
		"total":   len(articles),
		"success": success,
		"failed":  failed,
	}, "embedding 生成完成", c)
}

type VectorSearchRequest struct {
	Query string `form:"query" binding:"required"`
}

func (AiApi) SearchByVectorView(c *gin.Context) {
	var req VectorSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}

	vector, err := ai_service.GetEmbedding(req.Query)
	if err != nil {
		res.FailWithMsg("向量生成失败", c)
		return
	}

	script := elastic.NewScript(
		"cosineSimilarity(params.query_vector, 'content_vector') + 1.0",
	).Param("query_vector", vector)

	scoreQuery := elastic.NewScriptScoreQuery(
		elastic.NewTermQuery("status", 3),
		script,
	)

	result, err := global.ESClient.Search(models.ArticleModel{}.Index()).
		Query(scoreQuery).
		From(0).
		Size(10).
		Do(context.Background())
	if err != nil {
		res.FailWithMsg("检索失败", c)
		return
	}

	var list []map[string]interface{}
	for _, hit := range result.Hits.Hits {
		var article models.ArticleModel
		json.Unmarshal(hit.Source, &article)
		score := 0.0
		if hit.Score != nil {
			score = *hit.Score
		}
		list = append(list, map[string]interface{}{
			"id":    article.ID,
			"title": article.Title,
			"score": score,
		})
	}

	res.Success(list, "检索成功", c)
}
