package enum

type ArticleStatus int8

const (
	ArticleDraft     ArticleStatus = 0 // 草稿
	ArticleExamine   ArticleStatus = 1 // 审核中
	ArticlePublished ArticleStatus = 2 // 已发布
)
