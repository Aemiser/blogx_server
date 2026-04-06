package main

import (
	"blogx_server/api"
	"blogx_server/api/ai_api"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/conf"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models/enum"
	"blogx_server/service/ai_service"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// TestAuthMiddleware 测试用认证中间件（跳过 Redis 黑名单检查）
func TestAuthMiddleware(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Set("claims", claims)
}

// TestArticleAiView 直接测试 handler（跳过数据库查询，模拟 ES 返回数据）
func TestArticleAiView(c *gin.Context) {
	cr := middlerware.GetBind[ai_api.ArticleAiRequest](c)

	if !global.Config.Ai.Enable {
		res.SSEFail("站点未启用AI服务", c)
		return
	}

	// 模拟 ES 返回的文章数据（避免依赖真实数据库）
	content := `[{"title":"Go语言并发编程指南","abstract":"介绍goroutine和channel的使用","content":"详细内容","id":1},{"title":"Go Web开发实战","abstract":"基于Gin框架的Web开发","content":"详细内容","id":2}]`

	// 调用 AI 流式服务
	msgChan, err := ai_service.ChatStream(cr.Content, content)
	if err != nil {
		res.SSEFail("ai分析失败", c)
		return
	}
	for s := range msgChan {
		res.SSEOK(s, c)
	}
}

func main() {
	// 初始化配置
	global.Config = &conf.Config{
		Jwt: conf.Jwt{
			Secret: "test-secret-key-for-jwt-token-signing",
			Expire: 2400,
			Issuer: "taotao",
		},
		Ai: conf.Ai{
			Enable:    true,
			SecretKey: "sk-RGAlJrCPUw9tPqlefl6hVlwBJ2ljbfZ8AjSfu8pPi4b1FSX2",
		},
	}

	// 初始化错误码
	res.InitSysCode()

	// 生成有效 JWT token
	claims := jwts.Claims{
		UserID:   1,
		UserName: "testuser",
		Role:     enum.UserRole,
	}
	token, err := jwts.GetToken(claims)
	if err != nil {
		fmt.Printf("生成 token 失败: %v\n", err)
		return
	}
	fmt.Printf("生成 token: %s...\n\n", token[:50])

	// 创建测试路由
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	nr := r.Group("/api")
	app := api.App.AiApi

	// 路由 1: 完整链路（需要数据库，预期失败）
	nr.GET("ai/article", TestAuthMiddleware, middlerware.BindQueryMiddlerware[ai_api.ArticleAiRequest], app.ArticleAiView)

	// 路由 2: 模拟数据版本（绕过数据库）
	nr.GET("ai/article/mock", TestAuthMiddleware, middlerware.BindQueryMiddlerware[ai_api.ArticleAiRequest], TestArticleAiView)

	// 测试 1: 无 token 请求
	fmt.Println("=== 测试 1: 无 token 请求 ===")
	testWithoutToken(r)

	// 测试 2: 有效 token 但缺少 content 参数
	fmt.Println("\n=== 测试 2: 有效 token 但缺少 content 参数 ===")
	testWithTokenNoContent(r, token)

	// 测试 3: 无效 token
	fmt.Println("\n=== 测试 3: 无效 token ===")
	testWithInvalidToken(r)

	// 测试 4: 完整流程（模拟数据版本）
	fmt.Println("\n=== 测试 4: 完整流程（模拟数据版本）===")
	testWithTokenAndContentMock(r, token)

	// 测试 5: AI 禁用状态
	fmt.Println("\n=== 测试 5: AI 禁用状态 ===")
	testAiDisabled(token)

	// 测试 6: 完整流程（真实数据库版本 - 预期失败）
	fmt.Println("\n=== 测试 6: 完整流程（真实数据库版本 - 预期失败）===")
	testWithTokenAndContentReal(r, token)
}

func testWithoutToken(r *gin.Engine) {
	req := httptest.NewRequest(http.MethodGet, "/api/ai/article?content=Go语言", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("状态码: %d\n", w.Code)
	fmt.Printf("响应: %s\n", strings.TrimSpace(w.Body.String()))

	if w.Code == 200 && strings.Contains(w.Body.String(), "token") {
		fmt.Println("✓ 正确返回 token 相关错误")
	}
}

func testWithTokenNoContent(r *gin.Engine, token string) {
	req := httptest.NewRequest(http.MethodGet, "/api/ai/article", nil)
	req.Header.Set("token", token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("状态码: %d\n", w.Code)
	fmt.Printf("响应: %s\n", strings.TrimSpace(w.Body.String()))

	if w.Code == 200 && strings.Contains(w.Body.String(), "参数") {
		fmt.Println("✓ 正确返回参数错误")
	}
}

func testWithInvalidToken(r *gin.Engine) {
	req := httptest.NewRequest(http.MethodGet, "/api/ai/article?content=test", nil)
	req.Header.Set("token", "invalid-token-xxx")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("状态码: %d\n", w.Code)
	fmt.Printf("响应: %s\n", strings.TrimSpace(w.Body.String()))

	if w.Code == 200 && strings.Contains(w.Body.String(), "token") {
		fmt.Println("✓ 正确返回 token 无效错误")
	}
}

func testWithTokenAndContentMock(r *gin.Engine, token string) {
	start := time.Now()
	req := httptest.NewRequest(http.MethodGet, "/api/ai/article/mock?content=Go语言并发", nil)
	req.Header.Set("token", token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	elapsed := time.Since(start)
	fmt.Printf("状态码: %d\n", w.Code)

	body := w.Body.String()
	lines := strings.Split(body, "\n")
	fmt.Printf("响应行数: %d\n", len(lines))

	if len(lines) > 0 {
		fmt.Println("前 5 行响应:")
		limit := 5
		if len(lines) < limit {
			limit = len(lines)
		}
		for i := 0; i < limit; i++ {
			line := strings.TrimSpace(lines[i])
			if len(line) > 120 {
				fmt.Printf("  [%d] %s...\n", i+1, line[:120])
			} else if line != "" {
				fmt.Printf("  [%d] %s\n", i+1, line)
			}
		}
	}

	fmt.Printf("\n总耗时: %v\n", elapsed)

	if w.Code == 200 && strings.Contains(body, "data:") {
		fmt.Println("✓ SSE 流式响应成功")
	}
}

func testAiDisabled(token string) {
	// 临时禁用 AI
	global.Config.Ai.Enable = false
	defer func() { global.Config.Ai.Enable = true }()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	nr := r.Group("/api")
	nr.GET("ai/article/mock", TestAuthMiddleware, middlerware.BindQueryMiddlerware[ai_api.ArticleAiRequest], TestArticleAiView)

	req := httptest.NewRequest(http.MethodGet, "/api/ai/article/mock?content=test", nil)
	req.Header.Set("token", token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Printf("状态码: %d\n", w.Code)
	fmt.Printf("响应: %s\n", strings.TrimSpace(w.Body.String()))

	if w.Code == 200 && strings.Contains(w.Body.String(), "未启用") {
		fmt.Println("✓ 正确返回 AI 未启用错误")
	}
}

func testWithTokenAndContentReal(r *gin.Engine, token string) {
	start := time.Now()
	req := httptest.NewRequest(http.MethodGet, "/api/ai/article?content=Go语言", nil)
	req.Header.Set("token", token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	elapsed := time.Since(start)
	fmt.Printf("状态码: %d\n", w.Code)
	fmt.Printf("总耗时: %v\n", elapsed)

	if w.Code == 500 {
		fmt.Println("✓ 预期失败（无数据库连接）")
	}
}
