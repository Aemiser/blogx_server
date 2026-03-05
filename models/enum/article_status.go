package enum

type ArticleStatus int8

const (
	ArticleDraft     ArticleStatus = 1 // 草稿
	ArticleExamine   ArticleStatus = 2 // 审核中
	ArticlePublished ArticleStatus = 3 // 已发布
	ArticleFail      ArticleStatus = 4 // 审核失败
)
