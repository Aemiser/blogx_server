package ai_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/service/ai_service"
	"context"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type ArticleAiRequest struct {
	Content string `form:"content" binding:"required"`
}

func (AiApi) ArticleAiView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleAiRequest](c)

	if !global.Config.Ai.Enable {
		res.SSEFail("站点未启用AI服务", c)
		return
	}
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("status", 3)) // 必须是已发布状态
	// 使用 MatchPhraseQuery 进行短语匹配，避免分词导致的错误匹配
	query.Should(
		elastic.NewMatchQuery("title", cr.Content),
		elastic.NewMatchQuery("abstract", cr.Content),
		elastic.NewMatchQuery("content", cr.Content),
	)
	result, err := global.ESClient.Search(models.ArticleModel{}.Index()).
		Query(query).
		From(1).
		Size(10).
		Do(context.Background())
	if err != nil {
		source, _ := query.Source()
		byteData, _ := json.Marshal(source)
		logrus.Errorf("查询失败 %s \n %s", err, string(byteData))
		res.SSEOK("查询失败", c)
		return
	}

	var list []string
	for _, hit := range result.Hits.Hits {
		list = append(list, string(hit.Source))
	}
	content := "[" + strings.Join(list, ",") + "]"
	msgChan, err := ai_service.ChatStream(cr.Content, content)
	if err != nil {
		res.SSEFail("ai分析失败", c)
		return
	}
	for s := range msgChan {
		res.SSEOK(s, c)
	}

}
