package user_api

import "time"

type UserListResponse struct {
	ID            uint      `json:"id"`
	Nickname      string    `json:"nickname"`
	Username      string    `json:"username"`
	AvatarI       string    `json:"avatar"`
	ArticleCount  int       `json:"articleCount"`  // 发文数
	FansCount     int       `json:"fansCount"`     // 粉丝数
	FocusCount    int       `json:"focusCount"`    // 关注数
	IndexCount    int       `json:"indexCount"`    // 主页访问数
	CreatedAt     time.Time `json:"createdAt"`     // 注册时间
	LastLoginDate time.Time `json:"lastLoginDate"` // 最后登录时间
}
