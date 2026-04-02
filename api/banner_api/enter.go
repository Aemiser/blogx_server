package banner_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

type BannerApi struct {
}

type BannerCreateRequest struct {
	Cover string `json:"cover"`
	Href  string `json:"href"`
	Show  bool   `json:"show"`
}

func (BannerApi) BannerCreateView(c *gin.Context) {
	var req BannerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	// 入库
	err := global.Db.Create(&models.BannerModel{
		Cover: req.Cover,
		Href:  req.Href,
		Show:  req.Show,
	}).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.SuccessWithMsg("上传成功", c)
}
func (BannerApi) BannerRemoveView(c *gin.Context) {
	var req models.IDListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}

	// 查库是否存在
	var list []models.BannerModel
	global.Db.Find(&list, "id in  ?", req.IDList)

	if len(list) > 0 {
		global.Db.Delete(&list)
	}

	res.SuccessWithMsgf(c, "删除banner%d个，成功%d个", len(req.IDList), len(list))
	return

}

func (BannerApi) BannerUpdateView(c *gin.Context) {
	var cr models.IDRequest
	if err := c.ShouldBindUri(&cr); err != nil {
		res.FailWithError(err, c)
		return
	}
	var req BannerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	var model models.BannerModel
	err := global.Db.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	err = global.Db.Model(&model).Updates(map[string]any{
		"cover": req.Cover,
		"href":  req.Href,
		"show":  req.Show,
	}).Error

	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.SuccessWithMsg("修改成功", c)

}

type BannerListRequest struct {
	common.PageInfo
	Show bool `form:"show"`
}

func (BannerApi) BannerListView(c *gin.Context) {
	var req BannerListRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	list, count, _ := common.ListQuery(models.BannerModel{
		Show: req.Show,
	}, common.Options{
		PageInfo: req.PageInfo,
	})

	res.SuccessWithList(list, count, c)

}
