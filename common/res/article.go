package res

const (
	ArticleCollectExists     Code = 7001
	ArticleCollectCreateFail Code = 7002
	ArticleCollectNotFound   Code = 7003
	ArticleCollectUpdateFail Code = 7004
	ArticleCollectDeleteFail Code = 7005
	ArticleNotFound          Code = 7006
	ArticleCollectFull       Code = 7007
)

func InitArticleCode() {
	RegisterCode(ArticleCollectExists, "分类已存在")
	RegisterCode(ArticleCollectCreateFail, "创建收藏夹错误")
	RegisterCode(ArticleCollectNotFound, "分类不存在")
	RegisterCode(ArticleCollectUpdateFail, "更新收藏夹错误")
	RegisterCode(ArticleCollectDeleteFail, "删除收藏夹错误")
	RegisterCode(ArticleNotFound, "文章不存在")
	RegisterCode(ArticleCollectFull, "收藏夹已满")
}
