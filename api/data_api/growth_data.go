package data_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GrowthDataRequest struct {
	Type int8 `form:"type" binding:"required,oneof=1 2 3"` // 1 流量 2 文章发布量 3 用户注册量
}
type GrowthDataResponse struct {
	GrowthRate int      `json:"growthRate"`
	GrowthNum  int      `json:"growthNum"`
	DateList   []string `json:"dateList"`
	CountList  []int    `json:"countList"`
}

func (DataApi) GrowthDataView(c *gin.Context) {
	cr := middlerware.GetBind[GrowthDataRequest](c)

	type Table struct {
		Date  string `gorm:"column:date"`
		Count int    `gorm:"column:count"`
	}

	var dataList []Table

	now := time.Now()
	before7 := now.AddDate(0, 0, -6)

	switch cr.Type {
	case 1:
		global.Db.Model(models.SiteFlowModel{}).Where("created_at >=? and created_at <= ? ",
			before7.Format("2006-01-02")+" 00:00:00",
			now.Format("2006-01-02 15:04:05")).
			Select("DATE(created_at) as `date`", "sum(count) as `count`").
			Group("`date`").Scan(&dataList)
	case 2:
		global.Db.Model(models.ArticleModel{}).Where("created_at >=? and created_at <= ? and status = ?",
			before7.Format("2006-01-02")+" 00:00:00",
			now.Format("2006-01-02 15:04:05"),
			enum.ArticlePublished).
			Select("DATE(created_at) as `date`", "count(id) as `count`").
			Group("`date`").Scan(&dataList)
	case 3:
		global.Db.Model(models.UserModel{}).Where("created_at >=? and created_at <= ? ",
			before7.Format("2006-01-02")+" 00:00:00",
			now.Format("2006-01-02 15:04:05")).
			Select("DATE(created_at) as `date`", "count(id) as `count`").
			Group("`date`").Scan(&dataList)
	}

	var dateMap = map[string]int{}
	for _, model := range dataList {
		date := strings.Split(model.Date, "T")[0]
		dateMap[date] = model.Count
	}

	response := UserDataResponse{}
	for i := 0; i < 7; i++ {
		dateS := before7.AddDate(0, 0, i).Format("2006-01-02")
		count, _ := dateMap[dateS]
		response.DateList = append(response.DateList, dateS)
		response.CountList = append(response.CountList, count)
	}

	// 算增长
	response.GrowthNum = response.CountList[len(response.CountList)-1] - response.CountList[len(response.CountList)-2]
	if response.CountList[len(response.CountList)-2] == 0 {
		response.GrowthRate = 100
	}
	response.GrowthRate = int(float64(response.GrowthNum)/float64(response.CountList[len(response.CountList)-2])) * 100
	res.SuccessWithData(response, c)

}
