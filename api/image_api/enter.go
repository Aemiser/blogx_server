package image_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/log_service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ImageApi struct {
}

type ImageListRequest struct {
	models.ImageModel
	WebPath string `json:"webPath"`
}

func (ImageApi) ImageListView(c *gin.Context) {
	var req common.PageInfo
	if err := c.ShouldBindQuery(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	log := log_service.GetLog(c)
	log.SetLogType(enum.QueryLogType)
	log.SetTitle("<span style='color: #1890ff'>🖼️ 查看图片列表</span>")
	log.SetItem("📄 分页信息", fmt.Sprintf("第 <span style='color: #1890ff'>%d</span> 页，每页 <span style='color: #1890ff'>%d</span> 条", req.Page, req.Limit))
	_list, count, err := common.ListQuery(models.ImageModel{}, common.Options{
		PageInfo: req,
		Likes:    []string{"filename"},
		Debug:    true,
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	log.SetItem("📊 查询结果", fmt.Sprintf("共查询到 <span style='color: #52c41a; font-weight: bold'>%d</span> 张图片", count))

	var list = make([]ImageListRequest, 0)
	for _, model := range _list {
		list = append(list, ImageListRequest{
			ImageModel: model,
			WebPath:    model.WebPath(),
		})
	}
	if count > 0 {
		log.SetItem("图片预览", fmt.Sprintf("<div style='display: flex; gap: 8px; flex-wrap: wrap; max-height: 200px; overflow-y: auto'>%s</div>",
			func() string {
				var imgs string
				for _, img := range list {
					imgs += fmt.Sprintf("<img src='%s' style='max-width: 80px; max-height: 80px; object-fit: cover; border-radius: 4px;' title='%s'/>", img.WebPath, img.Filename)
				}
				return imgs
			}()))
	}
	res.SuccessWithList(list, count, c)

}

func (ImageApi) ImageRemoveView(c *gin.Context) {
	var req models.IDListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}
	log := log_service.GetLog(c)
	log.SetLogType(enum.OperationLogType)
	log.SetTitle("<span style='color: #ff4d4f'>🗑️ 删除图片</span>")
	log.SetItem("🗑️ 请求删除ID列表", fmt.Sprintf("<span style='color: #ff4d4f'>%v</span>", req.IDList))
	var list []models.ImageModel
	global.Db.Find(&list, "id in ?", req.IDList)

	log.SetItem("📋 查询到图片数", fmt.Sprintf("<span style='color: #1890ff'>%d</span>", len(list)))
	if len(list) > 0 {
		log.SetItem("🖼️ 图片预览", fmt.Sprintf("<div style='display: flex; gap: 8px; flex-wrap: wrap; max-height: 150px; overflow-y: auto'>%s</div>",
			func() string {
				var imgs string
				for _, img := range list {
					imgs += fmt.Sprintf("<img src='%s' style='max-width: 60px; max-height: 60px; object-fit: cover; border-radius: 4px;' title='%s'/>", img.WebPath, img.Filename)
				}
				return imgs
			}()))
	}

	var successCount, failCount int64
	if len(list) > 0 {
		successCount = global.Db.Delete(&list).RowsAffected
	}
	failCount = int64(len(list)) - successCount

	if successCount > 0 {
		log.SetItem("✅ 删除结果", fmt.Sprintf("<span style='color: #52c41a; font-weight: bold'>成功删除 %d 张图片</span>", successCount))
	} else {
		log.SetItem("❌ 删除结果", "<span style='color: #8c8c8c'>删除失败</span>")
	}
	if failCount > 0 {
		log.SetItem("⚠️ 失败数", fmt.Sprintf("<span style='color: #faad14'>%d</span>", failCount))
	}

	msg := fmt.Sprintf("操作成功，成功%d,失败%d", successCount, failCount)
	res.SuccessWithMsg(msg, c)

}
