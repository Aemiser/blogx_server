package main

import (
	"blogx_server/conf"
	"blogx_server/global"
	"blogx_server/service/ai_service"
	"fmt"
	"strings"
)

func main() {
	// 配置 AI 密钥
	global.Config = &conf.Config{
		Ai: conf.Ai{
			SecretKey: "sk-RGAlJrCPUw9tPqlefl6hVlwBJ2ljbfZ8AjSfu8pPi4b1FSX2",
		},
	}

	fmt.Println("=== 测试 ChatStream 流式调用 ===")
	fmt.Println()

	// 测试 1: 简单对话
	fmt.Println("--- 测试 1: 简单对话 ---")
	testSimpleChat()

	fmt.Println()

	// 测试 2: 带自定义 prompt
	fmt.Println("--- 测试 2: 带自定义 prompt ---")
	testWithPrompt()
}

func testSimpleChat() {
	msgChan, err := ai_service.ChatStream("请用一句话介绍你自己", "")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	var fullResponse strings.Builder
	for msg := range msgChan {
		fmt.Print(msg)
		fullResponse.WriteString(msg)
	}
	fmt.Println()
	fmt.Printf("\n完整回复长度: %d 字符\n", fullResponse.Len())
}

func testWithPrompt() {
	customPrompt := "你是一个专业的技术博主，请用简洁的语言回答："
	msgChan, err := ai_service.ChatStream("Go语言中goroutine和channel的关系是什么？", customPrompt)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	var fullResponse strings.Builder
	for msg := range msgChan {
		fmt.Print(msg)
		fullResponse.WriteString(msg)
	}
	fmt.Println()
	fmt.Printf("\n完整回复长度: %d 字符\n", fullResponse.Len())
}
