package ai_service

import (
	"blogx_server/global"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

//go:embed chat.prompt
var prompt string

const (
	bashurl = "https://api.chatanywhere.tech/v1/chat/completions"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Request struct {
	Model   string    `json:"model"`
	Message []Message `json:"messages"`
	Stream  bool      `json:"stream"`
}

type ChatResponse struct {
	Id      string `json:"id"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role        string        `json:"role"`
			Content     string        `json:"content"`
			Annotations []interface{} `json:"annotations"`
		} `json:"delta"`
		Message struct {
			Role        string        `json:"role"`
			Content     string        `json:"content"`
			Annotations []interface{} `json:"annotations"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Usage   struct {
		PromptTokens            int `json:"prompt_tokens"`
		CompletionTokens        int `json:"completion_tokens"`
		TotalTokens             int `json:"total_tokens"`
		CompletionTokensDetails struct {
			AudioTokens     int `json:"audio_tokens"`
			ReasoningTokens int `json:"reasoning_tokens"`
		} `json:"completion_tokens_details"`
		PromptTokensDetails struct {
			AudioTokens  int `json:"audio_tokens"`
			CachedTokens int `json:"cached_tokens"`
		} `json:"prompt_tokens_details"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

func bashRequest(r Request) (res *http.Response, err error) {
	method := "POST"
	byteData, _ := json.Marshal(r)
	req, err := http.NewRequest(method, bashurl, bytes.NewBuffer(byteData))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", global.Config.Ai.SecretKey))
	req.Header.Add("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)
	return
}

func Chat(content string) (msg string, err error) {
	r := Request{
		Model: "gpt-3.5-turbo",
		Message: []Message{
			{
				Role:    "system",
				Content: prompt,
			},
			{
				Role:    "user",
				Content: content,
			},
		},
		Stream: false,
	}
	res, err := bashRequest(r)
	if err != nil {
		logrus.Errorf("请求失败 %s", err)
		err = fmt.Errorf("AI 服务请求失败：%w", err)
		return
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	// 检查 HTTP 状态码
	if res.StatusCode != http.StatusOK {
		logrus.Errorf("AI 服务返回错误状态码：%d, 响应：%s", res.StatusCode, string(body))
		err = fmt.Errorf("AI 服务错误 (状态码 %d): 服务器繁忙，请稍后再试", res.StatusCode)
		return
	}

	//fmt.Println(string(body))
	var resp ChatResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logrus.Errorf("解析失败 %s  响应内容：%s", err, string(body))
		err = fmt.Errorf("AI 响应解析失败")
		return
	}

	if len(resp.Choices) == 0 {
		logrus.Warnf("未获取数据：%s", string(body))
		return
	}
	//fmt.Println(resp)
	msg = resp.Choices[0].Message.Content
	return
}
