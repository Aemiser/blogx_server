package article_api

import (
	"blogx_server/common"
	"blogx_server/common/jwts"
	"blogx_server/common/res"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"

	"github.com/gin-gonic/gin"
)

type ArticleListRequest struct {
	common.PageInfo
	Type       int8  `json:"type" binding:"required,oneof=1 2 3 "` // 1看别人的 2看自己的 3管理员看
	UserID     uint  `json:"userID"`
	CategoryID *uint `json:"categoryID"`
	Status     enum.ArticleStatus
}

type ArticleListResponse struct {
	models.ArticleModel
	UserTop  bool `json:"userTop"`  // 用户是否置顶
	AdminTop bool `json:"adminTop"` //管理员是否置顶
}

func (ArticleApi) ArticleListView(c *gin.Context) {
	cr := middlerware.GetBind[ArticleListRequest](c)
	switch cr.Type {
	case 1:
		// 查别人 用户ID 必填
		if cr.UserID == 0 {
			res.FailWithMsg("用户ID不能为空", c)
			return
		}
		// 查询受限
		if cr.Page > 2 || cr.Limit > 10 {
			res.FailWithMsg("查询更多，请登入", c)
			return
		}
		cr.Status = 0
	case 2:
		// 查自己
		claims, err := jwts.ParseTokenByGin(c)
		if err != nil {
			res.FailWithMsg("请登录", c)
			return
		}
		cr.UserID = claims.Claims.UserID
	case 3:
		// 管理员
		claims, err := jwts.ParseTokenByGin(c)
		if !(err == nil && claims.Claims.Role == enum.AdminRole) {
			res.FailWithMsg("角色错误", c)
			return
		}
	}
	_list, count, _ := common.ListQuery(models.ArticleModel{
		UserID:     cr.UserID,
		CategoryID: cr.CategoryID,
		Status:     cr.Status,
	}, common.Options{
		Likes:    []string{"title"},
		PageInfo: cr.PageInfo,
	})

	var list = make([]ArticleListResponse, 0)
	for _, model := range _list {
		model.Content = ""
		list = append(list, ArticleListResponse{
			ArticleModel: model,
		})
	}
	res.SuccessWithList(list, count, c)
}
