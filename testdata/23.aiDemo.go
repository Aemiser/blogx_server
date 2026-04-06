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
	msg, err := ai_service.ChatStream("给我推荐几个go文章", " [{\"abstract\":\"简介简介简介简介简介简介简介简介\",\"category_id\":0,\"collect_count\":58,\"comment_count\":0,\"content\":\"6.java后端文章2\",\"cover\":\"http://localhost:5175/uploads/images/617e6f7aae9361cccd3c1cf2597686fc.jpg\",\"created_at\":\"2026-03-25 20:12:38\",\"deleted_at\":null,\"digg_count\":12,\"id\":51,\"look_count\":345,\"open_comment\":1,\"status\":3,\"tag_list\":[\"java\"],\"title\":\"7.java后端文章2\",\"updated_at\":\"2026-03-27 20:12:39\",\"user_id\":5},{\"abstract\":\"简介简介简介简介简介简介简介简介\",\"category_id\":2,\"collect_count\":870708,\"comment_count\":0,\"content\":\"这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是正文这是 正文这是正文这是正文这是正文这是正文这是正文这是正文\",\"cover\":\"http://localhost:5175/uploads/images/617e6f7aae9361cccd3c1cf2597686fc.jpg\",\"created_at\":\"2026-03-30 15:56:23\",\"deleted_at\":null,\"digg_count\":234,\"id\":55,\"look_count\":345345,\"open_comment\":1,\"status\":3,\"tag_list\":[\"go\",\"后端\"],\"title\":\"go后端文章55\",\"updated_at\":\"2026-04-01 01:19:44\",\"user_id\":5},{\"abstract\":\"简介简介简介简介简介简介简介简介\",\"category_id\":null,\"collect_count\":707,\"comment_count\":0,\"content\":\"这是正文这是正文这是正文这是正文这是正文\\n\",\"cover\":\"http://localhost:5175/uploads/images/617e6f7aae9361cccd3c1cf2597686fc.jpg\",\"created_at\":\"2026-04-01 00:58:14\",\"deleted_at\":null,\"digg_count\":65,\"id\":57,\"look_count\":456,\"open_comment\":1,\"status\":3,\"tag_list\":[\"go\",\"后端\"],\"title\":\"go后端文章4\",\"updated_at\":\"2026-04-01 00:58:14\",\"user_id\":1},{\"abstract\":\"本文介绍MySQL主从同步的工作原理及其配置步骤，阐述如何通过此机制实现数据备份和负载均衡。\",\"category_id\":null,\"collect_count\":0,\"comment_count\":0,\"content\":\"# MySQL 主从同步技术概述\\n\\nMySQL 的主从同步是一种广泛应用于数据备份和负载均衡的策略。通过这一机制，数据可以高效地从一个主服务器（Master）复制到多 个从服务器（Slave）。这种架构不但提高了数据的可用性，还显著增强了系统的整体性能。\\n\\n## 主从同步的工作原理\\n\\n在 MySQL 中，主从同步主要依赖于**二进制日志**（Binary Log）和**复制线程**来实现。其具体工作步骤可概括为以下几个方面：\\n\\n1. **二进制日志的生成**  \\n   主服务器在执行任何写入操作时，会将这些变更记录在二进制日志中。该日志不仅记录了所有修改数据的 SQL 语句，还包括相关数据变更的信息，从而确保每一步操作的可追溯性。\\n\\n2. **从服务器读取日志**  \\n   从服务器定期与主服务器进行通信，以查询和获取更新的二进制日志。在此过程中，**I/O 线程**负责与主服务器的交互，确保从服务器获得最新的日志信息。\\n\\n3. **应用日志更新**  \\n   当从服务器成功获取到新日志后，它会利用 **SQL 线程**逐条执行这些日志中记载的 SQL 语句，从而对数据进行更新，确保从服务器的数据与主服务器保持同步。\\n\\n## 主从同步的配置步骤\\n\\n成功配置主从同步通常需要经过几个关键步骤，以确保每个环节均正确无误：\\n\\n1. **配置主服务器**  \\n   在主服务器的 配置文件（my.cnf）中，启用二进制日志功能，并添加以下参数：\\n   ```ini\\n   [mysqld]\\n   log-bin=mysql-bin\\n   server-id=1\\n   ```\\n\\n2. **创建复制用户**  \\n   在主服务器上，创建一个具有复制权限的用户，并分配相应的权限，如下所示 ：\\n   ```sql\\n   CREATE USER 'replicator'@'%' IDENTIFIED BY 'password';\\n   GRANT REPLICATION SLAVE ON *.* TO 'replicator'@'%';\\n   ```\\n\\n3. **获取主服务器状态**  \\n   接下来，执行以下命令以获取当前的二进制日志文件及其位置：\\n   ```sql\\n   SHOW MASTER STATUS;\\n   ```\\n\\n4. **配置从服务器**  \\n   在从服务器的配置文件中，设置一个唯一的服务器 ID，并输入主服务器的信息：\\n   ```ini\\n   [mysqld]\\n   server-id=2\\n   ```\\n\\n5. **建立与主服务器的连接**  \\n   使用以下命令在从服务器配置主服务器的连接信息：\\n   ```sql\\n   CHANGE MASTER TO\\n       MASTER_HOST='主服务器IP',\\n       MASTER_USER='replicator',\\n       MASTER_PASSWORD='password',\\n       MASTER_LOG_FILE='mysql-bin.000001',\\n       MASTER_LOG_POS=位置;\\n   ```\\n\\n6. **启动复制进程**  \\n   最后，在从服务器上启动复制线程：\\n   ```sql\\n   START SLAVE;\\n   ```\\n\\n## 监控主从同步状态\\n\\n为了确保同步过程的正常进行，可以通过执行以下命令来检查从服务器的同步状态：\\n```sql\\nSHOW SLAVE STATUS\\\\G\\n```\\n\\n在输出结果中，**Seconds_Behind_Master** 字段尤其值得关注，它表示从服务器与主服务器之间的时间延迟。这一指标有助于管理员及时发现和解决同步问题，从而确保系统的稳定运行。\",\"cover\":\"http://127.0.0.1:8080/uploads//d7fe66794e2deb9f9d76a4cdcfbe79e4.png\",\"created_at\":\"2026-04-04 17:13:00\",\"deleted_at\":null,\"digg_count\":0,\"id\":66,\"look_count\":0,\"open_comment\":1,\"status\":3,\"tag_list\":[\"MySQL\",\"主从同步\",\"数据备份\",\"负载均衡\",\"数据库管理\"],\"title\":\"MySQL 主从同步\",\"updated_at\":\"2026-04-04 17:13:33\",\"user_id\":5},{\"abstract\":\"本文介绍了一系列常用的 Linux 基础命令，涵盖文件与目录操作、文件操作、文件权限与所有权、系统信息与管理及网络相关命 令，旨在帮助新手顺利入门 Linux 环境。\",\"category_id\":null,\"collect_count\":0,\"comment_count\":0,\"content\":\"Linux 是一个强大的操作系统，掌握基础命令对新手至关重要。以下是一些常用的 Linux 基础命令，帮助您在 Linux 环境中顺利入门。\\n\\n## 1. 文件与目录操作\\n\\n- **`ls`**：列出目录内容。\\n  - 示例：`ls -la` 显示所有文件，包括隐藏文件，并提供详细信息。\\n\\n- **`cd`**：改变当前目录。\\n  - 示例：`cd /home/user/` 进入指定目录。\\n\\n- **`pwd`**：显示当前工作目录的完整路径。\\n\\n- **`mkdir`**：创建新目录。\\n  - 示例：`mkdir new_folder` 创建一个名为 `new_folder` 的新目录。\\n\\n- **`rmdir`**：删除空目录。\\n  - 示例：`rmdir old_folder` 删除名为 `old_folder` 的空目录。\\n\\n- **`rm`**：删除文件或 目录。\\n  - 示例：`rm file.txt` 删除 `file.txt` 文件；`rm -r folder/` 递归删除目录及其内容。\\n\\n## 2. 文件操作\\n\\n- **`cp`**：复制文件或目录。\\n  - 示例：`cp source.txt destination.txt` 将 `source.txt` 复制到 `destination.txt` 。\\n\\n- **`mv`**：移动或重命名文件或目录。\\n  - 示例：`mv oldname.txt newname.txt` 将 `oldname.txt` 重命名为 `newname.txt`。\\n\\n- **`touch`**：创建空文件或修改现有文件的时间戳。\\n  - 示例：`touch newfile.txt` 创建一个名为 `newfile.txt` 的空文件。\\n\\n- **`cat`**：查看文件内容。\\n  - 示例：`cat file.txt` 显示 `file.txt` 的内容。\\n\\n## 3. 文件权限与所有权\\n\\n- **`chmod`**：更改文件或目录的访问权限。\\n  - 示例：`chmod 755 file.sh` 设置 `file.sh` 的权限为可读、可写、可执行。\\n\\n- **`chown`**：更改文件或目录的所有者。\\n  - 示例：`chown user:group file.txt` 将 `file.txt` 的所有者更改为指定用户和组。\\n\\n## 4. 系统信息与管理\\n\\n- **`top`**：实时查看系统运行状态，包括 CPU 和内存占用。\\n\\n- **`df`**：查看文件系统的磁盘空间使用情况。\\n  - 示例：`df -h` 以人类可读的方式显示磁盘使用情况。\\n\\n- **`free`**：查看内存使用情况。\\n  - 示例：`free -m` 以 MB 为单位显示内存信息。\\n\\n- **`uname`**：显示系统信息。\\n  - 示例：`uname -a` 显示详细的系统和内核信息。\\n\\n## 5. 网络相关\\n\\n- **`ping`**：测试网络连接。\\n  - 示例：`ping google.com` 检查与 Google 服务器的连接。\\n\\n- **`ifconfig`**：查看或配置网络接口。\\n  - 示例：`ifconfig -a` 显示所有网络接口的详细信息。\\n\\n\",\"cover\":\"http://127.0.0.1:8080/uploads/images/6b3e448c40f403d0e319fc47005c9b21.webp\",\"created_at\":\"2026-04-03 12:08:51\",\"deleted_at\":null,\"digg_count\":0,\"id\":63,\"look_count\":0,\"open_comment\":1,\"status\":3,\"tag_list\":[\"Linux\",\"基础命令\",\"文件操作\",\"系统管理\",\"网络配置\"],\"title\":\"Linux 基础入门命令\",\"updated_at\":\"2026-04-03 12:08:51\",\"user_id\":1},{\"abstract\":\"简介简介简介简介简介简介简介简介\",\"category_id\":0,\"collect_count\":9,\"comment_count\":0,\"content\":\"#标题\\n\\njava正文\\n\",\"cover\":\"http://localhost:5175/uploads/images/617e6f7aae9361cccd3c1cf2597686fc.jpg\",\"created_at\":\"2026-03-26 20:57:15\",\"deleted_at\":\"2026-04-01 01:26:26\",\"digg_count\":45,\"id\":50,\"look_count\":23,\"open_comment\":0,\"status\":3,\"tag_list\":[\"java\",\"后端\"],\"title\":\"6.java后端文章\",\"updated_at\":\"2026-03-26 20:57:15\",\"user_id\":5},{\"abstract\":\"简介简介简介简介简介简介简介简介\",\"category_id\":2,\"collect_count\":98,\"comment_count\":0,\"content\":\"nihao hahahahha\\n\",\"cover\":\"http://localhost:5175/uploads/images/617e6f7aae9361cccd3c1cf2597686fc.jpg\",\"created_at\":\"2026-01-29 20:33:57\",\"deleted_at\":null,\"digg_count\":12,\"id\":52,\"look_count\":345,\"open_comment\":1,\"status\":3,\"tag_list\":[\"go\",\"后端2\"],\"title\":\"go后端文章2\",\"updated_at\":\"2026-04-01 01:02:33\",\"user_id\":5},{\"abstract\":\"简介简介简介简介简介简介简介简介\",\"category_id\":3,\"collect_count\":78070,\"comment_count\":0,\"content\":\"这是更新测试\\n\",\"cover\":\"http://127.0.0.1:8080/uploads//b9f488d9afc10da464160790b054cf70.png\",\"created_at\":\"2026-04-01 00:57:58\",\"deleted_at\":null,\"digg_count\":3453,\"id\":56,\"look_count\":345,\"open_comment\":1,\"status\":3,\"tag_list\":[\"go\",\"后端\"],\"title\":\"go后端文章4\",\"updated_at\":\"2026-04-05 01:43:45\",\"user_id\":5},{\"abstract\":\"这是一张新头像的展示，带有一副时尚而个性的风格。\",\"category_id\":null,\"collect_count\":0,\"comment_count\":0,\"content\":\"新 头像\\n![](\\nhttp://localhost:5175/uploads/images/u=2316432875,2297131572\\u0026fm=253\\u0026gp=0.jpg)\\n\\n\",\"cover\":\"http://127.0.0.1:8080/uploads/images/7553c5fdbeda33b84d0e771d4eddf0b2.jpg\",\"created_at\":\"2026-04-03 12:10:49\",\"deleted_at\":null,\"digg_count\":0,\"id\":64,\"look_count\":0,\"open_comment\":1,\"status\":3,\"tag_list\":[\"头像\",\"个人\",\"时尚\",\"展示\"],\"title\":\"新头像\",\"updated_at\":\"2026-04-03 12:10:49\",\"user_id\":1},{\"abstract\":\"本文介绍了 HTML 中的 div 标签，包括其基本使用、CSS 样式的应用以及语义化的重要性。\",\"category_id\":null,\"collect_count\":0,\"comment_count\":0,\"content\":\"### div 标签\\n\\n`<div>` 标签是 HTML 中的一个块级元素，它的全称是 &#34;division&#34;。该标签通常用于将网页内容分组，以便于布局和样式的管理。\\n\\n#### 基本使用\\n\\n使用 `<div>` 标签时，可以为其添加不同的属性，例如 `class` 和 `id`，以便通过 CSS 和 JavaScript 进行样式和交互的控制。\\n\\n```html\\n<div class=\\\"example-div\\\" id=\\\"unique-id\\\">\\n    这里是一个示例的 div 标签。\\n</div>\\n```\\n\\n#### CSS 样式\\n\\n通过 CSS，您可以轻松地为 `<div>` 标签定义样式。例如，可以设置背景颜色、边框、内外边距等属性：\\n\\n```css\\n.example-div {\\n    background-color: #f0f0f0;\\n    border: 1px solid #ccc;\\n    padding: 20px;\\n    margin: 10px 0;\\n}\\n```\\n\\n#### 语义化与使用场景\\n\\n虽然 `<div>` 标签非常有用，但由于它本身没有语义，因此在构建网页时，应该尽量使用其他具备语义意义的标签，如 `<header>`、`<footer>`、`<section>` 等。而 `<div>` 标签更多地用于需要特定样式或分组的场合。\\n\\n#### 小结\\n\\n在网页布局设计中，`<div>` 标签是一个不可或缺的工具，但应合理使用，以确保代码的可读性与结构的清晰性。通过与 CSS 和 JavaScript 的结合，能够打造出美观且功能丰富的网页元素。\\n\\n<iframe src=\\\"//player.bilibili.com/player.html?isOutside=true&amp;aid=116339171728674&amp;bvid=BV1Nn9TBhEYR&amp;cid=37202953877&amp;p=1\\\" scrolling=\\\"no\\\" border=\\\"0\\\" frameborder=\\\"no\\\" framespacing=\\\"0\\\" allowfullscreen=\\\"true\\\"></iframe>\",\"cover\":\"\",\"created_at\":\"2026-04-05 16:44:47\",\"deleted_at\":null,\"digg_count\":0,\"id\":74,\"look_count\":0,\"open_comment\":1,\"status\":3,\"tag_list\":[\"HTML\",\"div标签\",\"CSS\",\"网页布局\",\"前端开发\"],\"title\":\"深入了解 HTML 中的 div 标签\",\"updated_at\":\"2026-04-05 16:45:38\",\"user_id\":5}]")
	fmt.Println(err, msg)
	for item := range msg {
		fmt.Println(item)
	}
	//msg, err := ai_service.Chat("紧急！apifox被投毒，我已中招，赶紧自查\n发布时间：2026-03-26 （20 小时前）\napifox\n网络安全\n3月4日至3月22日 这期间打开过apifox的，赶紧看看自己有没有中招\n\nApifox 是一款 API 一体化协作平台，其桌面端应用基于 Electron 框架开发，提供 Windows、macOS、Linux 三平台客户端。因未严格启用 sandbox 参数，并暴露了 Node.js 的 API 接口，导致攻击者可通过 JS 控制 Apifox 的终端——三个平台均受影响\n\n简单来说：\n\napifox在打开过程中，会加载：\n\nhxxps://cdn[.]apifox[.]com/www/assets/js/apifox-app-event-tracking.min.js\n该文件正常大小为 34KB，但在 3 月 4 日之后可能会请求到被投毒的版本（77KB）。被投毒的 JS 文件会动态加载 hxxps://apifox[.]it[.]com/public/apifox-event.js（该域名非官方域名），在满足特定条件下加载攻击载荷，采集主机系统环境和敏感信息（SSH 密钥、Git 凭证、命令行历史、进程列表），上报到 hxxps://apifox[.]it[.]com/event/0/log。后续攻击者会控制主机拉取执行后门程序，并尝试发起横向攻击，控制更多有价值目标。\n\n如何自查\n\nwindows用户，去访问 %APPDATA%\\apifox\\Local Storage\\leveldb，去查看全部的二进制文件\n\n如果能看到 rl_mc 或 rl_headers ，那就说明中招了\n\n\n立即停用 Apifox 桌面端应用或者更新最新版本\n轮换所有 SSH 密钥（~/.ssh/ 下的全部密钥对）\n吊销所有 Git Personal Access Token（GitHub、GitLab 等）\n轮换 K8s 集群 OIDC Token 和 kubeconfig\n轮换 npm registry Token\n修改命令行历史中暴露的所有密码、Token 和 API Key\n审查服务器登录日志，检查是否有异常 SSH 登录 last命令\n原文链接：https://rce.moe/2026/03/25/apifox-supply-chain-attack-analysis")
	//fmt.Println(msg, err)

	//msgChan, err := ai_service.ChatStream("给我关于java的文章")
	//if err != nil {
	//	fmt.Println("错误:", err)
	//	return
	//}
	//
	//// 使用 for-range 循环接收流式数据
	//for s := range msgChan {
	//	fmt.Print(s) // 使用 Print 而不是 Println，避免换行
	//}
}
