package chat_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/ctype/chat_type"
	"blogx_server/models/enum/chat_msg_type"
	"blogx_server/models/enum/relationship_enum"
	"blogx_server/service/focus_service"
	"blogx_server/utils/xss"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var UP = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//	==> websocket地址 ==> websocket连接
//
// 用户id 	==> websocket地址 ==> websocket连接
//
//	==> websocket地址 ==> websocket连接
var OnlineMap = map[uint]map[string]*websocket.Conn{}

type ChatRequest struct {
	RevUserID uint                  `json:"revUserID"` // 发给谁
	MsgType   chat_msg_type.MsgType `json:"msgType"`   // 1 文本 2 图片 3 md
	Msg       chat_type.ChatMsg     `json:"msg"`       // 信息主体
}

type ChatResponse struct {
	ChatListResponse
}

func (ChatApi) ChatView(c *gin.Context) {
	claims, err := jwts.ParseTokenByGin(c)
	if err != nil || claims == nil {
		res.FailWithCode(res.SysUnauthorized, c)
		return
	}

	conn, err := UP.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorf("ws服务升级失败：%v", err)
		return
	}

	userID := claims.Claims.UserID
	var user models.UserModel
	err = global.Db.Take(&user, userID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}
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
	fmt.Println("服务开启:", OnlineMap)

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			if err == io.EOF {
				fmt.Println("客户端断开")
			}
			break
		}

		var req ChatRequest
		err2 := json.Unmarshal(p, &req)
		if err2 != nil {
			res.SendConnFailWithMsg(res.FailValueCode.String(), conn)
			continue

		}
		var revUser models.UserModel
		err1 := global.Db.Take(&revUser, req.RevUserID).Error
		if err1 != nil {
			res.SendConnFailWithMsg(res.ChatNotFound.String(), conn)
			continue
		}

		switch req.MsgType {
		case chat_msg_type.TextMsgType:
			if req.Msg.ContentMsg == nil || req.Msg.ContentMsg.Content == "" {
				res.SendConnFailWithMsg(res.ChatTextMsgEmpty.String(), conn)
				continue
			}
		case chat_msg_type.ImageMsgType:
			if req.Msg.ImagetMsg == nil || req.Msg.ImagetMsg.Src == "" {
				res.SendConnFailWithMsg(res.ChatImageMsgEmpty.String(), conn)
				continue
			}
		case chat_msg_type.MarkdownMsgType:
			if req.Msg.MarkdownMsg == nil || req.Msg.MarkdownMsg.Content == "" {
				res.SendConnFailWithMsg(res.ChatMarkdownMsgEmpty.String(), conn)
				continue
			}
			req.Msg.MarkdownMsg.Content = xss.Filter(req.Msg.MarkdownMsg.Content)
		default:
			res.SendConnFailWithMsg(res.ChatMsgTypeError.String(), conn)
			continue
		}

		relation := focus_service.CalcUserRelationship(userID, req.RevUserID)
		fmt.Printf("用户%d %d关系为：%d", userID, req.RevUserID, relation)
		switch relation {
		case relationship_enum.RelationStranger:
			var revUserMsgConf models.UserMessageConfModel
			err1 = global.Db.Take(&revUserMsgConf, "user_id= ?", revUser.ID).Error
			if err1 != nil {
				res.SendConnFailWithMsg(res.ChatPrivacyNotExist.String(), conn)
				continue
			}
			if !revUserMsgConf.OpenPrivateChat {
				res.SendConnFailWithMsg(res.ChatStrangerClosed.String(), conn)
				continue
			}

			var sendChatCount int64
			global.Db.Model(models.ChatModel{}).Where(" seed_user_id = ? and rev_user_id = ?",
				userID, req.RevUserID).Count(&sendChatCount)

			if sendChatCount >= 1 {
				res.SendConnFailInChatStrangerWithMsg(conn)
				continue
			}

		case relationship_enum.RelationFous, relationship_enum.RelationFans:
			var chatlist []models.ChatModel
			global.Db.Find(&chatlist, "date(created_at) = date (now()) and ((seed_user_id = ? and rev_user_id = ?) or (seed_user_id = ? and rev_user_id = ?))",
				userID, req.RevUserID, req.RevUserID, userID)

			var sendChatCount, revChatCount int
			for _, model := range chatlist {
				if model.SeedUserID == userID {
					sendChatCount++
				}

				if model.RevUserID == userID {
					revChatCount++
				}
			}
			if sendChatCount >= 1 && revChatCount == 0 {
				res.SendConnFailInChatLimitExceededWithMsg(conn)
				continue
			}
		}
		model := models.ChatModel{
			SeedUserID: userID,
			RevUserID:  req.RevUserID,
			MsgType:    req.MsgType,
			Msg:        req.Msg,
		}
		err = global.Db.Create(&model).Error
		if err != nil {
			res.SendConnFailWithMsg(res.ChatSendFailed.String(), conn)
			continue
		}
		item := ChatResponse{
			ChatListResponse: ChatListResponse{
				ChatModel:        model,
				SendUserNickname: user.Nickname,
				SendUserAvatar:   user.Avatar,
				RevUserNickname:  revUser.Nickname,
				RevUserAvatar:    revUser.Avatar,
			},
		}
		fmt.Println("发送给对方:", req.RevUserID)
		res.SendWsMsg(OnlineMap, req.RevUserID, item)
		item.IsMe = true
		fmt.Println("发送给自己:", item)
		res.SendConnOkWithData(item, conn)
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
	fmt.Println("服务关闭:", OnlineMap)
}

// canSendMessage 检查用户是否可以向对方发送消息
// 规则：如果当天已经发送过消息但对方未回复，则不能再次发送
func canSendMessage(userID, revUserID uint) bool {
	// 统计当天我发送的消息数量（我是发送者）
	var mySendCount int64
	global.Db.Model(&models.ChatModel{}).
		Where("date(created_at) = date(now()) AND seed_user_id = ? AND rev_user_id = ?",
			userID, revUserID).
		Count(&mySendCount)

	// 如果我还没发过消息，可以发送
	if mySendCount == 0 {
		return true
	}

	// 统计当天对方发送的消息数量（对方是发送者，我是接收者）
	var revSendCount int64
	global.Db.Model(&models.ChatModel{}).
		Where("date(created_at) = date(now()) AND seed_user_id = ? AND rev_user_id = ?",
			revUserID, userID).
		Count(&revSendCount)

	// 如果对方回复过，可以继续发送
	return revSendCount > 0
}
