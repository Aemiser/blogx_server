package xss

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//func Filter(content string) string {
//	contentDoc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(content)))
//	if err != nil {
//		return ""
//	}
//	contentDoc.Find("script").Remove()
//	contentDoc.Find("img").Remove()
//	contentDoc.Find("iframe").Remove()
//	return contentDoc.Text()
//}

// 生成随机占位符
func generatePlaceholder() string {
	b := make([]byte, 16)
	rand.Read(b)
	return "CODEBLOCK_" + hex.EncodeToString(b) + "_"
}

// Filter 过滤 HTML，但保护 Markdown 代码块
func Filter(content string) string {
	// 步骤1: 提取并保护代码块（包括 ``` 和 ` 两种格式）
	codeBlocks := make(map[string]string)

	// 保护 ``` 代码块（支持语言标识）
	tripleBacktickPattern := regexp.MustCompile("(?s)```[\\w]*\\n.*?```")
	content = tripleBacktickPattern.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := generatePlaceholder()
		codeBlocks[placeholder] = match
		return placeholder
	})

	// 保护 ` 行内代码
	inlineCodePattern := regexp.MustCompile("`[^`]+`")
	content = inlineCodePattern.ReplaceAllStringFunc(content, func(match string) string {
		placeholder := generatePlaceholder()
		codeBlocks[placeholder] = match
		return placeholder
	})

	// 步骤2: 处理 HTML 标签
	contentDoc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(content)))
	if err != nil {
		return content // 解析失败返回原文
	}

	// 移除危险标签
	contentDoc.Find("script").Remove()
	contentDoc.Find("img").Remove()

	// 处理 iframe 白名单
	contentDoc.Find("iframe").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists || src == "" {
			s.Remove()
			return
		}
		if !isInWhiteList(src) {
			s.Remove()
		}
		s.SetAttr("sandbox", "allow-scripts allow-same-origin allow-popups")
		s.SetAttr("referrerpolicy", "no-referrer")
	})

	// 获取处理后的 HTML
	body := contentDoc.Find("body")
	html, err := body.Html()
	if err != nil {
		return content
	}
	html = strings.TrimSpace(html)

	// 步骤3: 恢复代码块
	for placeholder, original := range codeBlocks {
		html = strings.ReplaceAll(html, placeholder, original)
	}

	return html
}

// 白名单检查（同之前）
func isInWhiteList(src string) bool {
	if strings.HasPrefix(src, "//") {
		src = "https:" + src
	}
	src = strings.ToLower(src)
	src = strings.TrimPrefix(src, "https://")
	src = strings.TrimPrefix(src, "http://")
	if idx := strings.Index(src, "/"); idx != -1 {
		src = src[:idx]
	}
	if idx := strings.Index(src, ":"); idx != -1 {
		src = src[:idx]
	}

	whiteList := map[string]bool{
		"player.bilibili.com": true,
	}
	return whiteList[src]
}
