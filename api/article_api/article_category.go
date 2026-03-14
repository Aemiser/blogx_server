package article_api

import (
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

type CategoryCreate struct {
	ID    uint   `json:"ID"`
	Title string `json:"Title" binding:"required,max=32"`
}

func (ArticleApi) CategoryCreateView(c *gin.Context) {
	cr := middlerware.GetBind[CategoryCreate](c)

	claims := jwts.GetClaimsByGin(c)
	var model models.CategoryModel
	if cr.ID == 0 {
		//创建
		err := global.Db.Take(&model, "user_id  = ? and title = ?", claims.Claims.UserID, cr.Title).Error
		if err == nil {
			res.FailWithMsg("分类已存在", c)
			return
		}

		err = global.Db.Create(&models.CategoryModel{
			Title:  cr.Title,
			UserID: claims.Claims.UserID,
		}).Error
		if err != nil {
			res.FailWithMsg("创建分类错误", c)
			return
		}
		res.SuccessWithMsg("创建分类成功", c)
		return
	}

	// 更新
	// 先根据 ID 查询出记录
	err := global.Db.First(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("分类不存在", c)
		return
	}

	// 再执行更新
	err = global.Db.Model(&model).Update("title", cr.Title).Error
	if err != nil {
		res.FailWithMsg("更新分类错误", c)
		return
	}
	res.SuccessWithMsg("更新分类成功", c)
}
