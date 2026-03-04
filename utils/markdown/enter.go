package markdown

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func MdToHtml(md string) string {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

func ExtractContent(content string, lenth int) (newcontent string, err error) {
	htmlcontent := MdToHtml(content)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlcontent)))
	if err != nil {
		return
	}
	htmlText := doc.Text()
	// 将字符串转换为 rune 切片来获取实际的字符数
	newcontent = htmlText
	runes := []rune(htmlText)
	if len(runes) > lenth {
		newcontent = string(runes[:lenth])
	}
	return
}
