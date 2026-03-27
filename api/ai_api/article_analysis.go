package ai_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/service/ai_service"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AIAnalysisRequest struct {
	Content string `json:"content" binding:"required"`
}

type AIAnalysisRespose struct {
	Title    string   `json:"title"`
	Abstract string   `json:"abstract"`
	Category string   `json:"category"`
	Tag      []string `json:"tag"`
}

func (AiApi) AIAnalysisView(c *gin.Context) {
	cr := middlerware.GetBind[AIAnalysisRequest](c)

	if !global.Config.Ai.Enable {
		res.FailWithMsg("站点未启用AI服务", c)
		return
	}

	msg, err := ai_service.Chat(cr.Content)
	if err != nil {
		logrus.Errorf("AI分析失败: %s %s", err, cr.Content)
		res.FailWithMsg("AI分析失败", c)
		return
	}
	var resp AIAnalysisRespose
	err = json.Unmarshal([]byte(msg), &resp)
	if err != nil {
		logrus.Errorf("解析失败: %s %s", err, msg)
		res.FailWithMsg("解析失败", c)
	}

	res.SuccessWithData(resp, c)
}
