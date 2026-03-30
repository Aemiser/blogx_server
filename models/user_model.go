package models

import (
	"blogx_server/models/enum"
	"math"
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	Model
	Username        string                  `gorm:"size:32" json:"username"`
	Nickname        string                  `gorm:"size:32" json:"nickname"`
	Avatar          string                  `gorm:"size:256" json:"avatar"`
	Abstract        string                  `gorm:"size:256" json:"abstract"`
	RegisterSource  enum.RegisterSourceType `json:"registerSource"` //注册来源
	CodeAge         int                     `json:"codeAge"`        //码龄
	Password        string                  `gorm:"size:64" json:"-"`
	Email           string                  `gorm:"size:256" json:"email"`
	OpenID          string                  `gorm:"size:64" json:"openID"` //第三方登入ID
	Role            enum.RoleType           `json:"role"`                  //角色:   1 管理员 2 普通用户 3 访客
	UserConfigModel *UserConfigModel        `gorm:"foreignKey:UserID" json:"-"`
	IP              string                  `gorm:"size:32" json:"ip"`
	Addr            string                  `gorm:"size:32" json:"addr"`
	ArticleList     []ArticleModel          `gorm:"foreignKey:UserID" json:"-"`
	LoginList       []UserLoginModel        `gorm:"foreignKey:UserID" json:"-"`
}

func (u UserModel) GetID() uint {
	return u.ID
}
func (u *UserModel) AfterCreate(tx *gorm.DB) (err error) {
	err = tx.Create(&UserConfigModel{UserID: u.ID,
		OpenCollect: true,
		OpenFollow:  true,
		OpenFans:    true,
		HomeStyleID: 1}).Error
	err = tx.Create(&UserMessageConfModel{
		UserID:             u.ID,
		OpenCommentMessage: true,
		OpenDiggMessage:    true,
		OpenPrivateChat:    true,
	}).Error
	return
}

func (u *UserModel) GetCodeAge() uint {
	sub := time.Now().Sub(u.CreatedAt)
	return uint(math.Ceil(sub.Hours() / 24 / 365))
}

type UserConfigModel struct {
	UserID             uint       `gorm:"primarykey;unique" json:"userID"`
	UserModel          UserModel  `gorm:"foreignKey:UserID" json:"-"`
	LikeTags           []string   `gorm:"type:longtext;serializer:json" json:"likeTags"`
	UpdataUsernameDate *time.Time `json:"updataUsernameDate"` // 上次修改用户名的时间
	OpenCollect        bool       `json:"openCollect"`        // 公开我的收藏
	OpenFollow         bool       `json:"openFollow"`         // 公开我的关注
	OpenFans           bool       `json:"openFans"`           // 公开我的粉丝
	HomeStyleID        uint       `json:"homeStyleID"`        // 主页样式ID
	LookCount          int        `json:"lookCount"`          // 主页的访问接口
}
