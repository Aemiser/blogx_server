package models

type UserMessageConfModel struct {
	UserID             uint      `gorm:"primarykey;unique" json:"userID"`
	UserModel          UserModel `gorm:"foreignKey:UserID" json:"-"`
	OpenCommentMessage bool      `json:"openCommentMessage"` //是否开启评论消息
	OpenDiggMessage    bool      `json:"openDiggMessage"`    //是否开启点赞
	OpenPrivateChat    bool      `json:"openPrivateChat"`    //是否开启私聊
}
