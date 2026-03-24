package chat_msg_type

type MsgType int8

const (
	TextMsgType MsgType = iota + 1
	ImageMsgType
	MarkdownMsgType

	MsgReadMsg = iota + 8
)
