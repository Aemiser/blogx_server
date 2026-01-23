package log_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

type LogApi struct {
}

type LogListRequest struct {
	Limit int    `form:"limit" json:"limit"`
	Page  int    `form:"page" json:"page"`
	Key   string `form:"key" json:"key"`
}

func (LogApi) LogListView(c *gin.Context) {
	// 分页查询 精确查询，模糊匹配
	var req LogListRequest
	err := c.ShouldBind(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	fmt.Println(req)

	var List []models.LogModel
	if req.Page >= 20 {
		req.Page = 1
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit == 0 || req.Limit > 100 {
		req.Limit = 10
	}
	offest := (req.Page - 1) * req.Limit
	global.Db.Debug().Offset(offest).Limit(req.Limit).Find(&List)

	var count int64
	global.Db.Debug().Model(&models.LogModel{}).Count(&count)

	res.FailWithList(List, int(count), c)
	return
}
