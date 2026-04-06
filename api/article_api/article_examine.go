package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/service/message_service"
	"blogx_server/service/text_service"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ArticleExamineRequest struct {
	ArticleID uint   `json:"articleID" binding:"required"`
	Status    uint   `json:"status" binding:"required"`
	Msg       string `json:"msg"`
}

func (ArticleApi) ArticleExamineView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleExamineRequest](c)

	var article models.ArticleModel
	err := global.Db.Take(&article, cr.ArticleID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	// 修改状态
	global.Db.Model(&article).Update("status", cr.Status)
	switch cr.Status {
	case 3: // 审核成功
		message_service.InsertSystemMessage(article.UserID, "管理员审核了你的文章", "文章审核成功", article.Title, fmt.Sprintf("/article/%d", article.ID))
		// 在全文列表中查找是否存在该文章的记录，如果不存在，说明在文章更新时被删除了，这里审核通过重新加入到里面
		textList := text_service.MdContentTransformation(article.ID, article.Title, article.Content)

		err = global.Db.Create(&textList).Error
		if err != nil {
			logrus.Errorf("创建文章ID为：%d的text失败：%v", article.ID, err)
		}
	case 4: // 审核失败
		message_service.InsertSystemMessage(article.UserID, "管理员审核了你的文章", fmt.Sprintf("文章审核失败，失败原因：%s", cr.Msg), "", "")
	}

	res.SuccessWithMsg("审核成功", c)

}
