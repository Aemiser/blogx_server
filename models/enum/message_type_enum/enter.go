package message_type_enum

type Type int8

const (
	CommentType Type = iota + 1
	ApplyType
	DiggARticleType
	UnDiggArticleType
	DiggCommentType
	UnDiggCommentType
	CollectArticleType
	UnCollectArticleType
	SystemType
)
