package data_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"time"

	"github.com/gin-gonic/gin"
)

type ArticleDataResponse struct {
	GrowthRate int      `json:"growthRate"`
	GrowthNum  int      `json:"growthNum"`
	DateList   []string `json:"dateList"`
	CountList  []int    `json:"countList"`
}

func (DataApi) ArticleDataView(c *gin.Context) {
	now := time.Now()

	before7 := now.AddDate(0, 0, -7)
	// 查询七天内的文章
	var articleList []models.ArticleModel
	global.Db.Find(&articleList, "created_at >=? and created_at <= ? and status = ?",
		before7.Format("2006-01-02")+" 00:00:00",
		now.Format("2006-01-02 15:04:05"),
		enum.ArticlePublished)

	var dateMap = map[string]int{}

	for _, model := range articleList {
		date := model.CreatedAt.Format("2006-01-02")
		count, ok := dateMap[date]
		if !ok {
			dateMap[date] = 1
			continue
		}
		dateMap[date] = count + 1
	}

	response := ArticleDataResponse{}
	for i := 0; i < 7; i++ {
		dateS := before7.AddDate(0, 0, i).Format("2006-01-02")
		count, _ := dateMap[dateS]
		response.DateList = append(response.DateList, dateS)
		response.CountList = append(response.CountList, count)
	}

	// 算增长
	response.GrowthNum = response.CountList[len(response.CountList)-1] - response.CountList[len(response.CountList)-2]
	response.GrowthRate = int(float64(response.GrowthNum)/float64(response.CountList[len(response.CountList)-2])) * 100
	res.SuccessWithData(response, c)
}
