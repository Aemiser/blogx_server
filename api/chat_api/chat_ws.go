package chat_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var UP = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//	==> websocket地址 ==> websocket连接
//
// 用户id 	==> websocket地址 ==> websocket连接
//
//	==> websocket地址 ==> websocket连接
var OnlineMap = map[uint]map[string]*websocket.Conn{}

func (ChatApi) ChatView(c *gin.Context) {
	// 用户认证
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil || claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}

	// 服务升级
	conn, err := UP.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorf("ws服务升级失败：%v", err)
		return
	}

	userID := claims.Claims.UserID
	addr := conn.RemoteAddr().String()
	addrMap, ok := OnlineMap[userID]
	if !ok { // 不存在即创建
		OnlineMap[userID] = map[string]*websocket.Conn{
			addr: conn,
		}

	} else { // 存在，但地址对应的连接不存在即追加
		_, ok1 := addrMap[addr]
		if !ok1 {
			OnlineMap[userID][addr] = conn
		}
	}
	fmt.Println("进入", OnlineMap)

	for {
		// 消息类型，消息，错误
		t, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err) // websocket: close 1005 (no status) :客户端断开
			if err == io.EOF {
				fmt.Println("客户端断开")
			}
			break
		}
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("你说的是：%s吗？", string(p))))
		fmt.Println(t, string(p))
	}
	defer conn.Close()

	addrMap2, ok2 := OnlineMap[userID]
	if ok2 {
		_, ok3 := addrMap2[addr]
		if ok3 {
			delete(OnlineMap[userID], addr)
		}
		if len(OnlineMap) == 0 {
			delete(OnlineMap, userID)
		}
	}
	fmt.Println("离开", OnlineMap)
	fmt.Println("服务关闭")
}
