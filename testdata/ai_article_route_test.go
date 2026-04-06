package main

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/conf"
	"blogx_server/global"
	"blogx_server/models/enum"
	"blogx_server/router"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

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
	fmt.Printf("生成 token: %s\n\n", token[:50]+"...")

	// 创建 Gin 引擎
	ginEngine := router.CreateTestEngine()

	// 测试 1: 无 token 请求（应该返回 401 未授权）
	fmt.Println("=== 测试 1: 无 token 请求 ===")
	testWithoutToken(ginEngine)

	// 测试 2: 有效 token 但缺少 content 参数（应该返回参数错误）
	fmt.Println("\n=== 测试 2: 有效 token 但缺少 content 参数 ===")
	testWithTokenNoContent(ginEngine, token)

	// 测试 3: 有效 token + 有效 content 参数（测试完整流程）
	fmt.Println("\n=== 测试 3: 有效 token + 有效 content 参数 ===")
	testWithTokenAndContent(ginEngine, token)
}

func testWithoutToken(engine *http.Handler) {
	req := httptest.NewRequest(http.MethodGet, "/api/ai/article?content=Go语言", nil)
	w := httptest.NewRecorder()
	(*engine).ServeHTTP(w, req)

	fmt.Printf("状态码: %d\n", w.Code)
	fmt.Printf("响应: %s\n", w.Body.String())

	if w.Code == 200 && strings.Contains(w.Body.String(), "token") {
		fmt.Println("✓ 正确返回 token 相关错误")
	}
}

func testWithTokenNoContent(engine *http.Handler, token string) {
	req := httptest.NewRequest(http.MethodGet, "/api/ai/article", nil)
	req.Header.Set("token", token)
	w := httptest.NewRecorder()
	(*engine).ServeHTTP(w, req)

	fmt.Printf("状态码: %d\n", w.Code)
	fmt.Printf("响应: %s\n", w.Body.String())

	if w.Code == 200 && strings.Contains(w.Body.String(), "参数") {
		fmt.Println("✓ 正确返回参数错误")
	}
}

func testWithTokenAndContent(engine *http.Handler, token string) {
	start := time.Now()
	req := httptest.NewRequest(http.MethodGet, "/api/ai/article?content=Go语言", nil)
	req.Header.Set("token", token)
	w := httptest.NewRecorder()
	(*engine).ServeHTTP(w, req)

	elapsed := time.Since(start)
	fmt.Printf("状态码: %d\n", w.Code)
	fmt.Printf("响应头 Content-Type: %s\n", w.Header().Get("Content-Type"))

	body := w.Body.String()
	// SSE 响应通常是多行的
	lines := strings.Split(body, "\n")
	fmt.Printf("响应行数: %d\n", len(lines))

	if len(lines) > 0 {
		fmt.Println("前 5 行响应:")
		limit := 5
		if len(lines) < limit {
			limit = len(lines)
		}
		for i := 0; i < limit; i++ {
			if len(lines[i]) > 100 {
				fmt.Printf("  [%d] %s...\n", i+1, lines[i][:100])
			} else {
				fmt.Printf("  [%d] %s\n", i+1, lines[i])
			}
		}
	}

	fmt.Printf("\n总耗时: %v\n", elapsed)

	if w.Code == 200 && strings.Contains(body, "data:") {
		fmt.Println("✓ SSE 流式响应成功")
	}
}
