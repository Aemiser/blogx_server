package message_type_enum

type Type int8

const (
	CommentType     Type = iota + 1 // 评论
	ApplyType                       // 回复
	DiggArticleType                 // 点赞文章
	UnDiggArticleType
	DiggCommentType // 点赞评论
	UnDiggCommentType
	CollectArticleType // 收藏文章
	UnCollectArticleType
	SystemType // 系统通知
)
