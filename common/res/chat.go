package res

const (
	ChatStranger         Code = 5001
	ChatNotFound         Code = 5002
	ChatMsgEmpty         Code = 5003
	ChatMsgTypeError     Code = 5004
	ChatPrivacyNotExist  Code = 5005
	ChatStrangerClosed   Code = 5006
	ChatLimitExceeded    Code = 5007
	ChatSendFailed       Code = 5008
	ChatTextMsgEmpty     Code = 5009
	ChatImageMsgEmpty    Code = 5010
	ChatMarkdownMsgEmpty Code = 5011
	ChatStrangerLimit    Code = 5012
	ChatDeleteFailed     Code = 5013
	ChatMsgDeleted       Code = 5014
	ChatSessionEmpty     Code = 5015
)

func InitChatCode() {
	RegisterCode(ChatStranger, "对方不是您的好友,只能发送一条消息")
	RegisterCode(ChatNotFound, "用户不存在")
	RegisterCode(ChatMsgEmpty, "消息内容为空")
	RegisterCode(ChatMsgTypeError, "消息类型错误")
	RegisterCode(ChatPrivacyNotExist, "隐私设置不存在")
	RegisterCode(ChatStrangerClosed, "对方未开启陌生人消息")
	RegisterCode(ChatLimitExceeded, "对方未回复的情况下，当天只能发送一条消息")
	RegisterCode(ChatSendFailed, "消息发送失败")
	RegisterCode(ChatTextMsgEmpty, "文本消息为空")
	RegisterCode(ChatImageMsgEmpty, "图片消息为空")
	RegisterCode(ChatMarkdownMsgEmpty, "markdown消息为空")
	RegisterCode(ChatStrangerLimit, "陌生人只能发送一条消息")
	RegisterCode(ChatDeleteFailed, "删除消息失败")
	RegisterCode(ChatMsgDeleted, "消息已删除")
	RegisterCode(ChatSessionEmpty, "聊天记录不存在")
}
