package message_type_enum

type Type int8

const(
	DiggArticleType Type = 1  // 点赞文章
	UnDiggArticleType Type =2 // 取消点赞文章
	DiggCommentType Type =3  // 点赞评论
	UnDiggCommentType Type =4   // 取消点赞评论
	CollectArticleType Type =5// 收藏文章
	UnCollectArticleType Type =6// 取消收藏文章
	SystemType Type =7



)