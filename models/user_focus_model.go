package models

import (
	"blogx_server/global"
	"blogx_server/models/enum/relationship_enum"
)

type UserFocusModel struct {
	Model
	UserID         uint      `json:"userID"`
	UserModel      UserModel `gorm:"foreignKey:UserID" json:"-"`
	FocusUserID    uint      `json:"focusID"`
	FocusUserModel UserModel `gorm:"foreignKey:FocusUserID" json:"-"`
}

func CalcUserRelationship(A, B uint) relationship_enum.Relation {
	var userFocusList []UserFocusModel
	global.Db.Find(&userFocusList,
		"(user_id = ? OR focus_user_id = ?) or (focus_user_id = ? OR  user_id= ?)",
		A, B, A, B)

	if len(userFocusList) == 2 {
		return relationship_enum.RelationFriends
	}

	if len(userFocusList) == 0 {
		return relationship_enum.RelationStranger
	}

	if userFocusList[0].FocusUserID == A {
		return relationship_enum.RelationFous
	}
	return relationship_enum.RelationFans
}
