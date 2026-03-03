package main

import (
	"blogx_server/utils/markdown"
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

var md = `
# 这是一级标题
> 这是引用
![xxx](adsfasdgsdfg)
`

func main() {
	html := markdown.MdToHtml(md)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(html)))
	if err != nil {
		fmt.Println(err)
		return
	}

	htmlText := doc.Text()
	//fmt.Println(htmlText)

	Abstract := htmlText[:5]
	fmt.Println(Abstract)
	Abstract = string([]rune(htmlText)[:5])
	fmt.Println(Abstract)

}
