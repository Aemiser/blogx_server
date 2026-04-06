package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/ai_service"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/olivere/elastic/v7"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.ESClient = core.EsConnect()

	if global.ESClient == nil {
		fmt.Println("ES 客户端未初始化，退出测试")
		return
	}

	fmt.Println("=== ES 检索功能测试 ===")
	fmt.Println()

	// 步骤 1: AI 提取关键词
	userInput := "给我推荐几个文章"
	fmt.Printf("用户输入: %s\n", userInput)

	msg, err := ai_service.KeywordChat(userInput)
	if err != nil {
		fmt.Printf("关键词提取失败: %v\n", err)
		return
	}
	fmt.Printf("AI 返回原始结果: [%s]\n", msg)
	fmt.Printf("AI 返回原始结果字节: %v\n", []byte(msg))

	var keywords []string
	err = json.Unmarshal([]byte(msg), &keywords)
	if err != nil {
		fmt.Printf("关键词解析失败: %v, 原始内容: %s\n", err, msg)
		return
	}
	fmt.Printf("关键词列表: %v\n", keywords)
	fmt.Println()

	// 步骤 2: 构建 ES 查询
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("status", 3))

	keywordQuery := elastic.NewBoolQuery()
	for _, keyword := range keywords {
		k := strings.ToLower(keyword)
		fmt.Printf("处理关键词: %s -> %s\n", keyword, k)

		keywordQuery.Should(
			elastic.NewMatchQuery("title", k),
			elastic.NewMatchQuery("abstract", k),
			elastic.NewMatchQuery("content", k),
		)
	}
	keywordQuery.MinimumNumberShouldMatch(1)
	query.Must(keywordQuery)

	// 打印 ES 查询 DSL
	source, _ := query.Source()
	byteData, _ := json.MarshalIndent(source, "", "  ")
	fmt.Printf("\nES 查询 DSL:\n%s\n", string(byteData))

	// 步骤 3: 执行查询
	result, err := global.ESClient.Search(models.ArticleModel{}.Index()).
		Query(query).
		From(0).
		Size(10).
		Do(context.Background())
	if err != nil {
		fmt.Printf("ES 查询失败: %v\n", err)
		return
	}

	fmt.Printf("\n总命中数: %d\n", result.TotalHits())
	fmt.Println()

	// 步骤 4: 检查结果
	if result.TotalHits() == 0 {
		fmt.Println("!!! 未找到任何匹配文章 !!!")
		fmt.Println("\n尝试直接查询 ES 确认数据是否存在...")
		testDirectQuery()
		return
	}

	foundDocker := false
	fmt.Println("匹配的文章:")
	for i, hit := range result.Hits.Hits {
		var article models.ArticleModel
		err = json.Unmarshal(hit.Source, &article)
		if err != nil {
			fmt.Printf("  [%d] 解析失败: %v\n", i+1, err)
			continue
		}
		title := article.Title
		hasDocker := strings.Contains(strings.ToLower(title), "docker") ||
			strings.Contains(strings.ToLower(article.Abstract), "docker") ||
			strings.Contains(strings.ToLower(article.Content), "docker")
		if hasDocker {
			foundDocker = true
		}
		fmt.Printf("  [%d] ID=%d, 标题=%s, 摘要=%s\n", i+1, article.ID, title, article.Abstract)
	}

	fmt.Println()
	if foundDocker {
		fmt.Println("✓ 结果中包含 Docker 相关文章")
	} else {
		fmt.Println("✗ 结果中未找到 Docker 相关文章")
	}
}

func testDirectQuery() {
	// 直接查询所有包含 docker 的文章（不带 status 过滤）
	query := elastic.NewBoolQuery().Should(
		elastic.NewMatchQuery("title", "docker"),
		elastic.NewMatchQuery("abstract", "docker"),
		elastic.NewMatchQuery("content", "docker"),
	)

	result, err := global.ESClient.Search(models.ArticleModel{}.Index()).
		Query(query).
		From(0).
		Size(5).
		Do(context.Background())
	if err != nil {
		fmt.Printf("直接查询失败: %v\n", err)
		return
	}

	fmt.Printf("直接查询 docker 命中数: %d\n", result.TotalHits())
	if result.TotalHits() > 0 {
		for i, hit := range result.Hits.Hits {
			var article models.ArticleModel
			json.Unmarshal(hit.Source, &article)
			fmt.Printf("  [%d] ID=%d, 标题=%s, status=%d\n", i+1, article.ID, article.Title, article.Status)
		}
	} else {
		fmt.Println("ES 中完全没有 docker 相关文章，可能是 River 同步问题")
	}
}
