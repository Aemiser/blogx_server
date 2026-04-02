package ai_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middlerware"
	"blogx_server/service/ai_service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ArticleImportRequest struct {
	Content string `json:"content" binding:"required"`
}

func (AiApi) ArticleImportView(c *gin.Context) {
	cr := middlerware.GetBind[AIAnalysisRequest](c)

	if !global.Config.Ai.Enable {
		res.FailWithMsg("站点未启用AI服务", c)
		return
	}

	msg, err := ai_service.ImportChat(cr.Content)
	if err != nil {
		logrus.Errorf("AI分析失败: %s %s", err, cr.Content)
		res.FailWithMsg("AI分析失败", c)
		return
	}
	res.SuccessWithData(msg, c)
}
