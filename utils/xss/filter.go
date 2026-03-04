package xss

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

func Filter(content string) string {
	contentDoc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(content)))
	if err != nil {
		return ""
	}
	contentDoc.Find("script").Remove()
	contentDoc.Find("img").Remove()
	contentDoc.Find("iframe").Remove()
	return contentDoc.Text()
}
