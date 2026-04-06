package main

import (
	"blogx_server/api"
	"blogx_server/api/article_api"
	"blogx_server/common/jwts"
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Db = core.InitDB()
	global.Redis = core.InitRedis()

	if global.Db == nil {
		fmt.Println("数据库未连接")
		return
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	nr := r.Group("/api")

	app := api.App.ArticleApi
	nr.GET("article", middlerware.BindQueryMiddlerware[article_api.ArticleListRequest], app.ArticleListView)

	// 查询一个有效用户生成 token
	var user models.UserModel
	global.Db.First(&user)
	if user.ID == 0 {
		fmt.Println("数据库中没有用户")
		return
	}
	fmt.Printf("使用用户: ID=%d, Username=%s\n", user.ID, user.Username)

	claims := jwts.Claims{
		UserID:   user.ID,
		UserName: user.Username,
		Role:     enum.UserRole,
	}
	token, _ := jwts.GetToken(claims)

	// 测试 1: collectID=15
	fmt.Println("\n=== 测试 1: collectID=15 ===")
	testArticleList(r, token, 15)

	// 测试 2: collectID=2
	fmt.Println("\n=== 测试 2: collectID=2 ===")
	testArticleList(r, token, 2)

	// 测试 3: 不带 collectID
	fmt.Println("\n=== 测试 3: 不带 collectID ===")
	testArticleListNoCollect(r, token)

	// 诊断: 检查收藏夹里的文章详情
	fmt.Println("\n=== 诊断: 收藏夹文章详情 ===")
	checkCollectArticles(15)
	checkCollectArticles(2)
}

func testArticleList(r *gin.Engine, token string, collectID uint) {
	url := fmt.Sprintf("/api/article?type=2&collectID=%d&page=1&limit=20", collectID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("token", token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("请求: %s\n", url)
	fmt.Printf("状态码: %d\n", w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if data, ok := resp["data"].(map[string]interface{}); ok {
		if list, ok := data["list"].([]interface{}); ok {
			fmt.Printf("返回文章数: %d\n", len(list))
			if count, ok := data["count"].(float64); ok {
				fmt.Printf("总数: %.0f\n", count)
			}
			for i, item := range list {
				if m, ok := item.(map[string]interface{}); ok {
					fmt.Printf("  [%d] ID=%.0f, 标题=%s\n", i+1, m["id"], m["title"])
				}
			}
		}
	} else {
		fmt.Printf("响应: %s\n", w.Body.String())
	}
}

func testArticleListNoCollect(r *gin.Engine, token string) {
	req := httptest.NewRequest(http.MethodGet, "/api/article?type=2&page=1&limit=20", nil)
	req.Header.Set("token", token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("请求: /api/article?type=2&page=1&limit=20\n")
	fmt.Printf("状态码: %d\n", w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if data, ok := resp["data"].(map[string]interface{}); ok {
		if list, ok := data["list"].([]interface{}); ok {
			fmt.Printf("返回文章数: %d\n", len(list))
			if count, ok := data["count"].(float64); ok {
				fmt.Printf("总数: %.0f\n", count)
			}
		}
	} else {
		fmt.Printf("响应: %s\n", w.Body.String())
	}
}
