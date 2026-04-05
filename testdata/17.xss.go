package main

import (
	"blogx_server/utils/xss"
	"fmt"
)

var md1 = `
# 这是一级标题
> 这是引用
![xxx](adsfasdgsdfg)

<script>alert(2) </script>
<img src="x" onerror="alert(2)" alt=""></img>
<iframe>alert(2)</iframe>
<iframe src="//player.bilibili.com/player.html?isOutside=true&aid=116339171728674&bvid=BV1Nn9TBhEYR&cid=37202953877&p=1" scrolling="no" border="0" frameborder="no" framespacing="0" allowfullscreen="true"></iframe>
`

func main() {
	//contentDoc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(md1)))
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//contentDoc.Find("script").Remove()
	//contentDoc.Find("img").Remove()
	//contentDoc.Find("iframe").Remove()
	//fmt.Println(contentDoc.Text())

	fmt.Println(xss.Filter(md1))
}
