package article_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"

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

type CategoryListRequest struct {
	common.PageInfo
	UserID uint `form:"UserID"`
	Type   int8 `form:"type" binding:"required,oneof=1 2 3 "` // 1查自己 2查别人 3后台
}
type CategoryListResponse struct {
	models.CategoryModel
	ArticleCount int `json:"ArticleCount"`
}

func (ArticleApi) CategoryListView(c *gin.Context) {
	cr := middlerware.GetBind[CategoryListRequest](c)

	switch cr.Type {
	case 1:
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			res.FailWithError(err, c)
			return
		}
		cr.UserID = claims.Claims.UserID
	case 2:

	case 3:
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			res.FailWithError(err, c)
			return
		}

		if claims.Claims.Role != enum.AdminRole {
			res.FailWithMsg("无权限", c)
			return
		}
	}

	_list, count, _ := common.ListQuery(models.CategoryModel{
		UserID: cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Preloads: []string{"ArticleList"},
		Likes:    []string{"title"}})

	var list = make([]CategoryListResponse, 0)
	for _, i2 := range _list {
		list = append(list, CategoryListResponse{
			CategoryModel: i2,
			ArticleCount:  len(i2.ArticleList),
		})
	}
	res.SuccessWithList(list, count, c)

}
func (ArticleApi) CategoryRemoveView(c *gin.Context) {
	var cr = middlerware.GetBind[models.IDListRequest](c)

	var list []models.CategoryModel

	query := global.Db.Where("id in ?", cr.IDList)
	claims := jwts.GetClaimsByGin(c)
	if claims.Claims.Role != enum.AdminRole {
		query.Where("user_id = ?", claims.Claims.UserID)
	}

	global.Db.Where(query).Find(&list)

	if len(list) > 0 {
		err := global.Db.Delete(&list).Error
		if err != nil {
			res.FailWithMsg("删除分类错误", c)
			return
		}
	}

	res.SuccessWithMsgf(c, "删除成功，成功删除分类 %d 条", len(list))

}

func (ArticleApi) CategoryOptionsView(c *gin.Context) {
	claims := jwts.GetClaimsByGin(c)

	var list []models.OptionsResponse[uint]
	global.Db.Model(&models.CategoryModel{}).Where("user_id = ?", claims.Claims.UserID).
		Select("id as value", "title as label").Scan(&list)

	res.SuccessWithData(list, c)
}
