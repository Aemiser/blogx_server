package ai_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"blogx_server/global"

	"github.com/sirupsen/logrus"
)

type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

func GetEmbedding(text string) (vector []float64, err error) {
	r := EmbeddingRequest{
		Model: "text-embedding-3-small",
		Input: []string{text},
	}
	byteData, _ := json.Marshal(r)
	req, err := http.NewRequest("POST", "https://api.chatanywhere.tech/v1/embeddings", bytes.NewBuffer(byteData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", global.Config.Ai.SecretKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embedding 请求失败: %w", err)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		logrus.Errorf("Embedding 服务返回错误状态码: %d, 响应: %s", res.StatusCode, string(body))
		return nil, fmt.Errorf("embedding 服务错误 (状态码 %d)", res.StatusCode)
	}

	var resp EmbeddingResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logrus.Errorf("Embedding 解析失败: %s, 响应: %s", err, string(body))
		return nil, fmt.Errorf("embedding 解析失败")
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("embedding 返回空数据")
	}

	return resp.Data[0].Embedding, nil
}
