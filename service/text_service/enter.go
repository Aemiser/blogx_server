package text_service

import (
	"blogx_server/models"
	"fmt"
	"strings"
)

func MdContentTransformation(model models.ArticleModel) (list []models.TextModel) {
	lines := strings.Split(model.Content, "\n")
	var headList []string
	var bodyList []string
	var body string
	var flag bool
	headList = append(headList, model.Title)
	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			flag = !flag
		}
		if !flag && strings.HasPrefix(line, "#") {
			// 标题行
			headList = append(headList, getHead(line))
			bodyList = append(bodyList, getBody(line))
			body = ""
			continue
		}
		body += line
	}

	if body != "" {
		bodyList = append(bodyList, getBody(body))
	}

	if len(headList) != len(bodyList) {
		fmt.Println("headList与bodyList长度不一致")
		fmt.Printf("%q  %d\n", headList, len(headList))
		fmt.Printf("%q  %d\n", bodyList, len(bodyList))
	}

	for i := 0; i < len(headList); i++ {
		list = append(list, models.TextModel{
			ArticleID: model.ID,
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
