package res

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

func SendConnFailWithMsg(msg string, conn *websocket.Conn) {
	data := Response{
		Code: FailValueCode,
		Data: empty,
		Msg:  msg,
	}
	byteData, _ := json.Marshal(data)
	conn.WriteMessage(websocket.TextMessage, byteData)
}

func SendConnOkWithData(data any, conn *websocket.Conn) {
	byteData, _ := json.Marshal(Response{
		Code: SuccessCode,
		Data: data,
		Msg:  "成功",
	})
	conn.WriteMessage(websocket.TextMessage, byteData)
}

func SendWsMsg(onlieMap map[uint]map[string]*websocket.Conn, userID uint, data any) {

	addrMap, ok := onlieMap[userID]
	// 没有这个接受人直接退出
	if !ok {
		return
	}
	// 编辑json格式信息
	byteData, _ := json.Marshal(Response{SuccessCode, data, "成功"})
	// 对每个客户端发送一条信息
	for _, conn := range addrMap {
		conn.WriteMessage(websocket.TextMessage, byteData)
	}
}
