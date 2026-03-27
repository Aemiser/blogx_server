package ai_service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Choices struct {
	Index int `json:"index"`
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
	Message struct {
	} `json:"message"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type StreamDate struct {
	Id                string    `json:"id"`
	Choices           []Choices `json:"choices"`
	Created           int       `json:"created"`
	Model             string    `json:"model"`
	Object            string    `json:"object"`
	SystemFingerprint string    `json:"system_fingerprint"`
}

func ChatStream(content string) (msgChan chan string, err error) {
	msgChan = make(chan string, 10) // 添加缓冲避免阻塞
	r := Request{
		Model: "gpt-3.5-turbo",
		Message: []Message{
			{
				Role:    "system",
				Content: "你是一名叫熊大的客服人工智能助手",
			},
			{
				Role:    "user",
				Content: content,
			},
		},
		Stream: true,
	}
	res, err := bashRequest(r)
	if err != nil {
		logrus.Errorf("请求失败 %s", err)
		err = fmt.Errorf("AI 服务请求失败：%w", err)
		close(msgChan)
		return
	}

	// 检查 HTTP 状态码
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		logrus.Errorf("AI 服务返回错误状态码：%d, 响应：%s", res.StatusCode, string(body))
		err = fmt.Errorf("AI 服务错误 (状态码 %d): 服务器繁忙，请稍后再试", res.StatusCode)
		close(msgChan)
		return
	}

	scanner := bufio.NewScanner(res.Body)
	scanner.Split(bufio.ScanLines)
	go func() {
		defer res.Body.Close() // 在 goroutine 结束时关闭
		defer close(msgChan)
		for scanner.Scan() {
			text := scanner.Text()
			if text == "" {
				continue
			}

			// 检查是否以 "data: " 开头
			if len(text) < 6 || text[:6] != "data: " {
				continue
			}

			data := text[6:]
			if data == "[DONE]" {
				return
			}
			var item StreamDate
			err = json.Unmarshal([]byte(data), &item)
			if err != nil {
				logrus.Errorf("解析失败 %s %s", err, data)
				continue
			}
			if len(item.Choices) == 0 {
				continue
			}
			content := item.Choices[0].Delta.Content
			if content != "" {
				msgChan <- content
			}
		}
	}()
	return
}
