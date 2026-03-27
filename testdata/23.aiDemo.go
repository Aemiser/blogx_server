package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/ai_service"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func chatv1() {

	url := "https://api.chatanywhere.tech/v1/chat/completions"
	method := "POST"

	payload := strings.NewReader(`{
    "model": "gpt-3.5-turbo",
    "messages": [
      {
        "role": "system",
        "content": "你是一名叫帅涛的人工只能助手"
      },
      {
        "role": "user",
        "content": "你是谁？"
      }
    ]
  }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", global.Config.Ai.SecretKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}

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

func streamchatv1() {
	url := "https://api.chatanywhere.tech/v1/chat/completions"
	method := "POST"

	payload := strings.NewReader(`{
    "model": "gpt-3.5-turbo",
    "messages": [
      {
        "role": "system",
        "content": "你是一名叫帅涛的人工只能助手"
      },
      {
        "role": "user",
        "content": "你是谁？"
      }
    ],
	"stream":true
  }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", global.Config.Ai.SecretKey))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}

		data := text[6:]
		if data == "[DONE]" {
			return
		}
		//fmt.Println(text)
		var item StreamDate
		err = json.Unmarshal([]byte(data), &item)
		if err != nil {
			logrus.Errorf("解析失败 %s %s", err, data)
			continue
		}
		if len(item.Choices) == 0 {
			continue
		}
		fmt.Println(item.Choices[0].Delta.Content)
	}
}
func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.InitLogrus()
	global.Db = core.InitDB()

	//chatv1()
	//streamchatv1()

	//msg, err := ai_service.Chat("紧急！apifox被投毒，我已中招，赶紧自查\n发布时间：2026-03-26 （20 小时前）\napifox\n网络安全\n3月4日至3月22日 这期间打开过apifox的，赶紧看看自己有没有中招\n\nApifox 是一款 API 一体化协作平台，其桌面端应用基于 Electron 框架开发，提供 Windows、macOS、Linux 三平台客户端。因未严格启用 sandbox 参数，并暴露了 Node.js 的 API 接口，导致攻击者可通过 JS 控制 Apifox 的终端——三个平台均受影响\n\n简单来说：\n\napifox在打开过程中，会加载：\n\nhxxps://cdn[.]apifox[.]com/www/assets/js/apifox-app-event-tracking.min.js\n该文件正常大小为 34KB，但在 3 月 4 日之后可能会请求到被投毒的版本（77KB）。被投毒的 JS 文件会动态加载 hxxps://apifox[.]it[.]com/public/apifox-event.js（该域名非官方域名），在满足特定条件下加载攻击载荷，采集主机系统环境和敏感信息（SSH 密钥、Git 凭证、命令行历史、进程列表），上报到 hxxps://apifox[.]it[.]com/event/0/log。后续攻击者会控制主机拉取执行后门程序，并尝试发起横向攻击，控制更多有价值目标。\n\n如何自查\n\nwindows用户，去访问 %APPDATA%\\apifox\\Local Storage\\leveldb，去查看全部的二进制文件\n\n如果能看到 rl_mc 或 rl_headers ，那就说明中招了\n\n\n立即停用 Apifox 桌面端应用或者更新最新版本\n轮换所有 SSH 密钥（~/.ssh/ 下的全部密钥对）\n吊销所有 Git Personal Access Token（GitHub、GitLab 等）\n轮换 K8s 集群 OIDC Token 和 kubeconfig\n轮换 npm registry Token\n修改命令行历史中暴露的所有密码、Token 和 API Key\n审查服务器登录日志，检查是否有异常 SSH 登录 last命令\n原文链接：https://rce.moe/2026/03/25/apifox-supply-chain-attack-analysis")
	//fmt.Println(msg, err)

	msgChan, err := ai_service.ChatStream("给我关于java的文章")
	if err != nil {
		fmt.Println("错误:", err)
		return
	}

	// 使用 for-range 循环接收流式数据
	for s := range msgChan {
		fmt.Print(s) // 使用 Print 而不是 Println，避免换行
	}
}
