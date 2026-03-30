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

type CollectCreate struct {
	ID       uint   `json:"ID"`
	Title    string `json:"Title" binding:"required,max=32"`
	Cover    string `json:"cover" `
	Abstract string `json:"abstract" binding:"required,max=256"`
}

func (ArticleApi) CollectCreateView(c *gin.Context) {
	cr := middlerware.GetBind[CollectCreate](c)

	claims := jwts.GetClaimsByGin(c)
	var model models.CollectModel
	if cr.ID == 0 {
		//创建
		err := global.Db.Take(&model, "user_id  = ? and title = ?", claims.Claims.UserID, cr.Title).Error
		if err == nil {
			res.FailWithMsg("分类已存在", c)
			return
		}

		err = global.Db.Create(&models.CollectModel{
			Title:    cr.Title,
			UserID:   claims.Claims.UserID,
			Cover:    cr.Cover,
			Abstract: cr.Abstract,
		}).Error
		if err != nil {
			res.FailWithMsg("创建收藏夹错误", c)
			return
		}
		res.SuccessWithMsg("创建收藏夹成功", c)
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
	err = global.Db.Model(&model).Updates(map[string]any{
		"title":    cr.Title,
		"cover":    cr.Cover,
		"abstract": cr.Abstract,
	}).Error
	if err != nil {
		res.FailWithMsg("更新收藏夹错误", c)
		return
	}
	res.SuccessWithMsg("更新收藏夹成功", c)
}

type CollectListRequest struct {
	common.PageInfo
	UserID    uint `form:"UserID"`
	Type      int8 `form:"type" binding:"required,oneof=1 2 3 "` // 1查自己 2查别人 3后台
	ArticleID uint `form:"articleID"`
}
type CollectListResponse struct {
	models.CollectModel
	ArticleCount int    `json:"ArticleCount"`
	Nickname     string `json:"nickname,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	ArticlceUse  bool   `json:"articlceUse"`
}

func (ArticleApi) CollectListView(c *gin.Context) {
	cr := middlerware.GetBind[CollectListRequest](c)

	var preloads = []string{"ArticleList"}
	switch cr.Type {
	case 1:
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			res.FailWithError(err, c)
			return
		}
		cr.UserID = claims.Claims.UserID
	case 2:
		var userconf models.UserConfigModel
		err := global.Db.Take(&userconf, "user_id = ?", cr.UserID).Error
		if err != nil {
			res.FailWithMsg("用户不存在", c)
			return
		}

		if !userconf.OpenCollect {
			res.FailWithMsg("用户未开放收藏功能", c)
			return
		}
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
		preloads = append(preloads, "UserModel")
	}

	_list, count, _ := common.ListQuery(models.CollectModel{
		UserID: cr.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"title"},
		Preloads: preloads})

	var list = make([]CollectListResponse, 0)
	for _, i2 := range _list {
		item := CollectListResponse{
			CollectModel: i2,
			ArticleCount: len(i2.ArticleList),
			Nickname:     i2.UserModel.Nickname,
			Avatar:       i2.UserModel.Avatar}
		for _, model := range i2.ArticleList {
			if model.ID == cr.ArticleID {
				item.ArticlceUse = true
				break
			}
		}
		list = append(list)
	}
	res.SuccessWithList(list, count, c)

}
func (ArticleApi) CollectRemoveView(c *gin.Context) {
	var cr = middlerware.GetBind[models.IDListRequest](c)

	var list []models.CollectModel

	query := global.Db.Where("id in ?", cr.IDList)
	claims := jwts.GetClaimsByGin(c)
	if claims.Claims.Role != enum.AdminRole {
		query.Where("user_id = ?", claims.Claims.UserID)
	}

	global.Db.Where(query).Find(&list)

	if len(list) > 0 {
		err := global.Db.Delete(&list).Error
		if err != nil {
			res.FailWithMsg("删除收藏夹错误", c)
			return
		}
	}

	res.SuccessWithMsgf(c, "删除成功，成功删除分类 %d 条", len(list))

}
