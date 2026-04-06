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

	"github.com/olivere/elastic/v7"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.ESClient = core.EsConnect()

	if global.ESClient == nil {
		fmt.Println("ES 客户端未初始化")
		return
	}

	fmt.Println("=== ES 向量检索 RAG v2 测试 ===")
	fmt.Println()

	fmt.Println("--- 步骤 1: 更新 ES mapping 添加 content_vector 字段 ---")
	addVectorMapping()

	fmt.Println("\n--- 步骤 2: 测试 embedding 生成 ---")
	testEmbedding()

	fmt.Println("\n--- 步骤 3: 为所有文章生成 embedding ---")
	generateAllEmbeddings()

	fmt.Println("\n--- 步骤 4: 测试向量检索 (docker) ---")
	testVectorSearch("给我推荐几个docker文章")

	fmt.Println("\n--- 步骤 5: 测试模糊语义检索 (容器化部署) ---")
	testVectorSearch("容器化部署")

	fmt.Println("\n--- 步骤 6: 测试通用推荐 ---")
	testVectorSearch("给我推荐几篇文章")
}

func addVectorMapping() {
	body := map[string]interface{}{
		"properties": map[string]interface{}{
			"content_vector": map[string]interface{}{
				"type":       "dense_vector",
				"dims":       1536,
				"index":      true,
				"similarity": "cosine",
			},
		},
	}

	_, err := global.ESClient.PutMapping().
		Index(models.ArticleModel{}.Index()).
		BodyJson(body).
		Do(context.Background())
	if err != nil {
		fmt.Printf("添加 mapping 失败 (可能已存在): %v\n", err)
	} else {
		fmt.Println("✓ content_vector 字段添加成功")
	}
}

func testEmbedding() {
	vector, err := ai_service.GetEmbedding("docker 容器技术")
	if err != nil {
		fmt.Printf("embedding 生成失败: %v\n", err)
		return
	}
	fmt.Printf("✓ embedding 生成成功, 维度: %d\n", len(vector))
	fmt.Printf("  前5个值: %.4f, %.4f, %.4f, %.4f, %.4f\n",
		vector[0], vector[1], vector[2], vector[3], vector[4])
}

func generateAllEmbeddings() {
	result, err := global.ESClient.Search(models.ArticleModel{}.Index()).
		Query(elastic.NewTermQuery("status", 3)).
		Size(100).
		Do(context.Background())
	if err != nil {
		fmt.Printf("查询文章失败: %v\n", err)
		return
	}

	fmt.Printf("共找到 %d 篇已发布文章\n", result.TotalHits())

	success := 0
	failed := 0

	for i, hit := range result.Hits.Hits {
		var article models.ArticleModel
		json.Unmarshal(hit.Source, &article)

		text := article.Title + " " + article.Abstract
		if text == " " {
			continue
		}

		vector, err := ai_service.GetEmbedding(text)
		if err != nil {
			fmt.Printf("  [%d] 文章 %d embedding 失败: %v\n", i+1, article.ID, err)
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
			fmt.Printf("  [%d] 文章 %d ES 更新失败: %v\n", i+1, article.ID, err)
			failed++
			continue
		}

		success++
		fmt.Printf("  [%d/%d] 文章 %d (%s) 完成\n", i+1, result.TotalHits(), article.ID, article.Title)
	}

	fmt.Printf("\n✓ 完成: 成功 %d, 失败 %d\n", success, failed)
}

func testVectorSearch(query string) {
	fmt.Printf("用户输入: %s\n", query)

	vector, err := ai_service.GetEmbedding(query)
	if err != nil {
		fmt.Printf("  embedding 生成失败: %v\n", err)
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
		Size(5).
		Do(context.Background())
	if err != nil {
		fmt.Printf("  向量检索失败: %v\n", err)
		fmt.Println("  尝试调试：打印查询 DSL...")
		src, _ := scoreQuery.Source()
		dsl, _ := json.MarshalIndent(src, "", "  ")
		fmt.Printf("  DSL: %s\n", string(dsl))
		return
	}

	fmt.Printf("  命中: %d 篇\n", result.TotalHits())
	if result.TotalHits() > 0 {
		for i, hit := range result.Hits.Hits {
			var article models.ArticleModel
			json.Unmarshal(hit.Source, &article)
			score := 0.0
			if hit.Score != nil {
				score = *hit.Score
			}
			fmt.Printf("  [%d] ID=%d, 标题=%s, 分数=%.4f\n",
				i+1, article.ID, article.Title, score)
		}
	} else {
		fmt.Println("  未找到匹配文章")
	}
}
