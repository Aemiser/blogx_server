package main

import (
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

var md1 = `
# 这是一级标题
> 这是引用
![xxx](adsfasdgsdfg)

<script>alert(2) </script>
<img src="x" onerror="alert(2)" alt=""></img>
<iframe>alert(2)</iframe>
`

func main() {
	contentDoc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(md1)))
	if err != nil {
		fmt.Println(err)
		return
	}
	contentDoc.Find("script").Remove()
	contentDoc.Find("img").Remove()
	contentDoc.Find("iframe").Remove()
	fmt.Println(contentDoc.Text())
}
