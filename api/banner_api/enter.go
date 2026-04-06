package banner_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/log_service"
	"fmt"

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
	log := log_service.GetLog(c)
	log.SetTitle("<span style='color: #1890ff'>➕ 创建Banner</span>")
	log.SetItem("封面图片", fmt.Sprintf("<img src='%s' style='max-width: 200px; border-radius: 4px;'/>", req.Cover))
	log.SetItem("跳转链接", fmt.Sprintf("<a href='%s' target='_blank' style='color: #1890ff'>%s</a>", req.Href, req.Href))
	showText := "❌ 不显示"
	if req.Show {
		showText = "✅ 显示"
	}
	log.SetItem("显示状态", fmt.Sprintf("<span style='font-weight: bold'>%s</span>", showText))

	// 入库
	err := global.Db.Create(&models.BannerModel{
		Cover: req.Cover,
		Href:  req.Href,
		Show:  req.Show,
	}).Error
	if err != nil {
		log.SetItemError("创建失败", err)
		res.FailWithError(err, c)
		return
	}
	log.SetItem("创建结果", "<span style='color: #52c41a; font-weight: bold'>✅ 创建成功</span>")
	res.SuccessWithMsg("上传成功", c)
}
func (BannerApi) BannerRemoveView(c *gin.Context) {
	var req models.IDListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}

	log := log_service.GetLog(c)
	log.SetTitle("<span style='color: #ff4d4f'>🗑️ 删除Banner</span>")
	log.SetItem("请求删除ID列表", fmt.Sprintf("<span style='color: #ff4d4f'>%v</span>", req.IDList))

	// 查库是否存在
	var list []models.BannerModel
	global.Db.Find(&list, "id in  ?", req.IDList)

	log.SetItem("查询到Banner数", fmt.Sprintf("<span style='color: #1890ff'>%d</span>", len(list)))
	if len(list) > 0 {
		var covers []string
		for _, model := range list {
			covers = append(covers, model.Cover)
		}
		log.SetItem("Banner封面", fmt.Sprintf("<div style='display: flex; gap: 8px; flex-wrap: wrap'>%s</div>",
			func() string {
				var imgs string
				for _, cover := range covers {
					imgs += fmt.Sprintf("<img src='%s' style='max-width: 100px; border-radius: 4px;'/>", cover)
				}
				return imgs
			}()))

		global.Db.Delete(&list)
		log.SetItem("删除结果", fmt.Sprintf("<span style='color: #52c41a; font-weight: bold'>✅ 成功删除 %d 个Banner</span>", len(list)))
	} else {
		log.SetItem("删除结果", "<span style='color: #8c8c8c'>未找到任何Banner</span>")
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

	log := log_service.GetLog(c)
	log.SetTitle("<span style='color: #faad14'>✏️ 更新Banner</span>")
	log.SetItem("Banner ID", fmt.Sprintf("<span style='color: #1890ff'>%d</span>", cr.ID))

	log.SetItem("📋 修改前", fmt.Sprintf("<div style='background: #fff7e6; padding: 8px; border-radius: 4px; border-left: 3px solid #faad14'>封面: <img src='%s' style='max-width: 80px;'/><br>链接: %s<br>显示: %v</div>", model.Cover, model.Href, model.Show))

	log.SetItem("📝 修改后", fmt.Sprintf("<div style='background: #e6f7ff; padding: 8px; border-radius: 4px; border-left: 3px solid #1890ff'>封面: <img src='%s' style='max-width: 80px;'/><br>链接: %s<br>显示: %v</div>", req.Cover, req.Href, req.Show))

	err = global.Db.Model(&model).Updates(map[string]any{
		"cover": req.Cover,
		"href":  req.Href,
		"show":  req.Show,
	}).Error

	if err != nil {
		log.SetItemError("更新失败", err)
		res.FailWithError(err, c)
		return
	}
	log.SetItem("更新结果", "<span style='color: #52c41a; font-weight: bold'>✅ 更新成功</span>")
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
