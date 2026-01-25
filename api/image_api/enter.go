package image_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
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
	_list, count, err := common.ListQuery(models.ImageModel{}, common.Options{
		PageInfo: req,
		Likes:    []string{"filename"},
		Debug:    true,
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	var list = make([]ImageListRequest, 0)
	for _, model := range _list {
		list = append(list, ImageListRequest{
			ImageModel: model,
			WebPath:    model.WebPath(),
		})
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
	log.ShowResponse()
	log.ShowRequest()
	var list []models.ImageModel
	global.Db.Find(&list, "id in ?", req.IDList)

	var successCount, failCount int64
	if len(list) > 0 {
		successCount = global.Db.Delete(&list).RowsAffected
	}
	failCount = int64(len(list)) - successCount
	msg := fmt.Sprintf("操作成功，成功%d,失败%d", successCount, failCount)
	res.SuccessWithMsg(msg, c)

}
