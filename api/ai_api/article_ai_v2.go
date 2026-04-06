package ai_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/service/ai_service"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

func (AiApi) ArticleAiViewV2(c *gin.Context) {
	cr := middlerware.GetBind[ArticleAiRequest](c)

	if !global.Config.Ai.Enable {
		res.SSEFail("站点未启用AI服务", c)
		return
	}

	vector, err := ai_service.GetEmbedding(cr.Content)
	if err != nil {
		res.SSEFail("向量生成失败", c)
		return
	}
	fmt.Printf("生成向量维度: %d\n", len(vector))

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
		logrus.Errorf("向量检索失败: %s", err)
		res.SSEFail("检索失败", c)
		return
	}

	fmt.Printf("向量检索命中: %d 篇文章\n", result.TotalHits())

	var list []string
	for _, hit := range result.Hits.Hits {
		var article models.ArticleModel
		err = json.Unmarshal(hit.Source, &article)
		if err != nil {
			logrus.Error("json.Unmarshal err:", err)
			continue
		}
		item := ArticleAIProRequest{
			ID:       article.ID,
			Title:    article.Title,
			Abstract: article.Abstract,
		}
		bytedata, err1 := json.Marshal(item)
		if err1 != nil {
			logrus.Error("json.Marshal err:", err1)
			continue
		}
		list = append(list, string(bytedata))
	}

	content := "[" + strings.Join(list, ",") + "]"
	fmt.Println("拼接的提示词：", content)

	msgChan, err := ai_service.ChatStream(cr.Content, content)
	if err != nil {
		res.SSEFail("ai分析失败", c)
		return
	}
	for s := range msgChan {
		res.SSEOK(s, c)
	}
}
