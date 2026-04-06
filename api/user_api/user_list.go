package user_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/middlerware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/log_service"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type UserListRequest struct {
	common.PageInfo
}
type UserListResponse struct {
	ID           uint   `json:"id"`
	Nickname     string `json:"nickname"`
	Username     string `json:"username"`
	AvatarI      string `json:"avatar"`
	IP           string `json:"ip"`
	Addr         string `json:"addr"`
	ArticleCount int    `json:"articleCount"` // 发文数
	//FansCount     int       `json:"fansCount"`     // 粉丝数
	//FocusCount    int       `json:"focusCount"`    // 关注数
	IndexCount    int           `json:"indexCount"`    // 主页访问数
	CreatedAt     time.Time     `json:"createdAt"`     // 注册时间
	LastLoginDate time.Time     `json:"lastLoginDate"` // 最后登录时间
	Role          enum.RoleType `json:"role"`
}

func (UserApi) UserListView(c *gin.Context) {
	cr := middlerware.GetBind[UserListRequest](c)

	log := log_service.GetLog(c)
	log.SetLogType(enum.QueryLogType)
	log.SetTitle("<span style='color: #1890ff'>👥 查看用户列表</span>")
	log.SetItem("📄 分页信息", fmt.Sprintf("第 <span style='color: #1890ff'>%d</span> 页，每页 <span style='color: #1890ff'>%d</span> 条", cr.Page, cr.Limit))

	_list, count, _ := common.ListQuery(models.UserModel{}, common.Options{
		Likes:    []string{"nickname", "username"},
		Preloads: []string{"ArticleList", "LoginList"},
		PageInfo: cr.PageInfo,
	})

	var list = make([]UserListResponse, 0)
	for _, model := range _list {
		item := UserListResponse{
			ID:           model.ID,
			Nickname:     model.Nickname,
			Username:     model.Username,
			AvatarI:      model.Avatar,
			IP:           model.IP,
			Addr:         model.Addr,
			ArticleCount: len(model.ArticleList),
			IndexCount:   1000,
			Role:         model.Role,
		}

		if len(model.LoginList) > 0 {
			item.LastLoginDate = model.LoginList[len(model.LoginList)-1].CreatedAt
		}
		list = append(list, item)
	}

	log.SetItem("📊 查询结果", fmt.Sprintf("共查询到 <span style='color: #52c41a; font-weight: bold'>%d</span> 位用户", count))
	if count > 0 {
		log.SetItem("用户详情", fmt.Sprintf("<div style='background: #f5f5f5; padding: 8px; border-radius: 4px; max-height: 200px; overflow-y: auto'>%s</div>",
			func() string {
				var result string
				for i, u := range list {
					roleIcon := "👤"
					if u.Role == enum.AdminRole {
						roleIcon = "👑"
					}
					result += fmt.Sprintf("%d. %s %s (ID: %d, 文章: %d, IP: %s)<br>", i+1, roleIcon, u.Nickname, u.ID, u.ArticleCount, u.IP)
				}
				return result
			}()))
	}
	res.SuccessWithList(list, count, c)
	return
}
