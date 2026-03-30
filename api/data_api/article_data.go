package data_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type ArticleYearDataResponse struct {
	GrowthRate int      `json:"growthRate"`
	GrowthNum  int      `json:"growthNum"`
	DateList   []string `json:"dateList"`
	CountList  []int    `json:"countList"`
}

func (DataApi) ArticleYearDataView(c *gin.Context) {
	now := time.Now()
	// 获取 12 个月前那个月的 1 号，避免月末溢出问题
	before12 := now.AddDate(0, -12, 0)
	startDate := time.Date(before12.Year(), before12.Month(), 1, 0, 0, 0, 0, before12.Location())
	fmt.Println("起始日期:", startDate.Format("2006-01-02"))
	var dataList []Table
	global.Db.Model(models.ArticleModel{}).Where("created_at >=? and created_at <= ? and status = ?",
		before12.Format("2006-01-02")+" 00:00:00",
		now.Format("2006-01-02 15:04:05"),
		enum.ArticlePublished).
		Select("MONTH(created_at) as `date`", "count(id) as `count`").
		Group("`date`").Scan(&dataList)
	var dateMap = map[string]int{}
	for _, model := range dataList {
		date := model.Date
		dateMap[date] = model.Count
	}

	response := ArticleYearDataResponse{}
	for i := 0; i < 12; i++ {
		date := startDate.AddDate(0, i+1, 0)
		count, _ := dateMap[date.Format("1")]
		response.DateList = append(response.DateList, date.Format("2006-01"))
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
