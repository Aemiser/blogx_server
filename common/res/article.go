package res

const (
	ArticleCollectExists      Code = 7001
	ArticleCollectCreateFail  Code = 7002
	ArticleCollectNotFound    Code = 7003
	ArticleCollectUpdateFail  Code = 7004
	ArticleCollectDeleteFail  Code = 7005
	ArticleNotFound           Code = 7006
	ArticleCollectFull        Code = 7007
	ArticleCreateFail         Code = 7008
	ArticleUpdateFail         Code = 7009
	ArticleDeleteFail         Code = 7010
	ArticleDiggFail           Code = 7011
	ArticleLookFail           Code = 7012
	ArticleExamineFail        Code = 7013
	ArticleNotOpenComment     Code = 7014
	ArticleTagNotFound        Code = 7015
	ArticleCategoryNotFound   Code = 7016
	ArticleCategoryCreateFail Code = 7017
	ArticleCategoryUpdateFail Code = 7018
	ArticleCategoryDeleteFail Code = 7019
)

func InitArticleCode() {
	RegisterCode(ArticleCollectExists, "分类已存在")
	RegisterCode(ArticleCollectCreateFail, "创建收藏夹错误")
	RegisterCode(ArticleCollectNotFound, "分类不存在")
	RegisterCode(ArticleCollectUpdateFail, "更新收藏夹错误")
	RegisterCode(ArticleCollectDeleteFail, "删除收藏夹错误")
	RegisterCode(ArticleNotFound, "文章不存在")
	RegisterCode(ArticleCollectFull, "收藏夹已满")
	RegisterCode(ArticleCreateFail, "文章创建失败")
	RegisterCode(ArticleUpdateFail, "文章更新失败")
	RegisterCode(ArticleDeleteFail, "文章删除失败")
	RegisterCode(ArticleDiggFail, "点赞失败")
	RegisterCode(ArticleLookFail, "浏览记录失败")
	RegisterCode(ArticleExamineFail, "文章审核失败")
	RegisterCode(ArticleNotOpenComment, "文章未开放评论")
	RegisterCode(ArticleTagNotFound, "标签不存在")
	RegisterCode(ArticleCategoryNotFound, "分类不存在")
	RegisterCode(ArticleCategoryCreateFail, "创建分类错误")
	RegisterCode(ArticleCategoryUpdateFail, "更新分类错误")
	RegisterCode(ArticleCategoryDeleteFail, "删除分类错误")
}
