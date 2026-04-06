package ai_api

import (
	"blogx_server/common"
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

type ArticleAiRequest struct {
	Content string `form:"content" binding:"required"`
}

type ArticleAIProRequest struct {
	ID       uint   `json:"id"`
	Abstract string `json:"abstract"`
	Title    string `json:"title"`
}

func (AiApi) ArticleAiView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleAiRequest](c)

	if !global.Config.Ai.Enable {
		res.SSEFail("站点未启用AI服务", c)
		return
	}

	var content string
	// 服务降级
	if global.ESClient == nil {
		list, _, _ := common.ListQuery(models.ArticleModel{}, common.Options{
			Likes: []string{"title", "abstract"},
			PageInfo: common.PageInfo{
				Limit: 10,
				Page:  1,
			},
		})

		byteData, _ := json.Marshal(list)
		content = string(byteData)
	} else {
		// 提取用户语句里的技术关键词
		msg, err := ai_service.KeywordChat(cr.Content)
		if err != nil {
			res.SSEFail("ai分析失败", c)
			return
		}
		fmt.Println(msg)
		var keywords []string
		err = json.Unmarshal([]byte(msg), &keywords)
		if err != nil {
			logrus.Errorf("解析失败: %s %s", err, msg)
			res.SSEFail("解析失败", c)
			return
		}
		fmt.Println("关键词列表：", keywords)
		query := elastic.NewBoolQuery()
		query.Must(elastic.NewTermQuery("status", 3)) // 必须是已发布状态

		if len(keywords) > 0 {
			keywordQuery := elastic.NewBoolQuery()
			for _, keyword := range keywords {
				k := strings.ToLower(keyword) // 转小写

				fmt.Printf("处理关键词: %s -> %s\n", keyword, k)

				keywordQuery.Should(
					elastic.NewMatchQuery("title", k),
					elastic.NewMatchQuery("abstract", k),
					elastic.NewMatchQuery("content", k),
				)
			}
			keywordQuery.MinimumNumberShouldMatch(1)
			query.Must(keywordQuery)
		}
		result, err := global.ESClient.Search(models.ArticleModel{}.Index()).
			Query(query).
			From(0).
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
			//fmt.Println("文章：", article)
		}
		content = "[" + strings.Join(list, ",") + "]"
	}
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
