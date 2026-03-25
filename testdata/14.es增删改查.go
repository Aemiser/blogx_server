package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models"
	"context"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
)

func DocCreate() {
	user := models.ArticleModel{
		Model: models.Model{gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
		}},
		Title:   "涛涛知道",
		Content: "这是内容",
		UserID:  1,
	}
	indexResponse, err := global.ESClient.Index().Index(user.Index()).BodyJson(user).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", indexResponse)
}

// 批量添加
func DocCreateBatch() {

	list := []models.ArticleModel{
		{
			Model: models.Model{gorm.Model{
				ID:        2,
				CreatedAt: time.Now(),
			}},
			Title:   "涛涛知道1",
			Content: "这是内容",
			UserID:  2,
		},
		{
			Model: models.Model{gorm.Model{
				ID:        3,
				CreatedAt: time.Now(),
			}},
			Title:   "涛涛知道2",
			Content: "这是内容",
			UserID:  3,
		},
	}

	bulk := global.ESClient.Bulk().Index(models.ArticleModel{}.Index()).Refresh("true")
	for _, model := range list {
		req := elastic.NewBulkCreateRequest().Doc(model)
		bulk.Add(req)
	}
	res, err := bulk.Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Succeeded())
}

// 批量查询
func DocFind() {

	limit := 2
	page := 1
	from := (page - 1) * limit
	query := elastic.NewBoolQuery()
	res, err := global.ESClient.Search(models.ArticleModel{}.Index()).Query(query).From(from).Size(limit).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	count := res.Hits.TotalHits.Value // 总数
	fmt.Println(count)
	for _, hit := range res.Hits.Hits {
		fmt.Println(string(hit.Source))
	}
}

// 删除文档
func DocDelete() {
	deleteResponse, err := global.ESClient.Delete().
		Index(models.ArticleModel{}.Index()).Id("sJAwV5wBXSZSyP8i9nm-").Refresh("true").Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(deleteResponse)
}

// 批量删除文档
func DocDeleteBatch() {
	idList := []string{
		"tGcofYkBWS69Op6QHJ2g",
	}
	bulk := global.ESClient.Bulk().Index(models.ArticleModel{}.Index()).Refresh("true")
	for _, s := range idList {
		req := elastic.NewBulkDeleteRequest().Id(s)
		bulk.Add(req)
	}
	res, err := bulk.Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Succeeded()) // 实际删除的文档切片
}

func updata() {
	global.ESClient.Update().Index(models.ArticleModel{}.Index()).Id("sZA8V5wBXSZSyP8i13mX").Doc(map[string]any{
		"title": "这是修改后的标题",
	}).Do(context.Background())
}
func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.ESClient = core.EsConnect()

	DocCreate()
	//updata()
	//DocFind( )
	//DocDelete()
}
