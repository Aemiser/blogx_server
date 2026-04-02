package text_service

import (
	"strings"
)

type TextModel struct {
	ArticleID uint   `json:"article_id"`
	Head      string `json:"head"`
	Body      string `json:"body"`
}

func MdContentTransformation(id uint, title, content string) (list []TextModel) {
	lines := strings.Split(content, "\n")
	var headList []string
	var bodyList []string
	var body string
	var flag bool
	// 先添加文章标题
	headList = append(headList, title)
	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			flag = !flag
		}
		if !flag && strings.HasPrefix(line, "#") {
			// 遇到 Markdown 标题：保存之前累积的正文，然后添加新标题
			bodyList = append(bodyList, getBody(body))
			headList = append(headList, getHead(line))
			body = ""
			continue
		}
		body += line
	}

	// 最后一段正文添加到 bodyList，如果为空则添加空字符串
	bodyList = append(bodyList, getBody(body))

	// 确保两个列表长度一致
	for i := 0; i < len(headList); i++ {
		list = append(list, TextModel{
			ArticleID: id,
			Head:      headList[i],
			Body:      bodyList[i],
		})
	}
	return
}
func getHead(head string) string {
	s := strings.TrimSpace(strings.Join(strings.Split(head, " ")[1:], " "))
	return s
}
func getBody(body string) string {
	return strings.TrimSpace(body)
}
