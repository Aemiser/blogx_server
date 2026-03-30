package focus_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum/relationship_enum"
	"fmt"
)

func CalcUserRelationship(A, B uint) relationship_enum.Relation {
	var userFocusList []models.UserFocusModel
	global.Db.Find(&userFocusList,
		"(user_id = ? and focus_user_id = ?) or (focus_user_id = ? and  user_id= ?)",
		A, B, A, B)
	fmt.Println(userFocusList)
	fmt.Println(len(userFocusList))
	if len(userFocusList) == 2 {
		return relationship_enum.RelationFriends
	}

	if len(userFocusList) == 0 {
		return relationship_enum.RelationStranger
	}

	if userFocusList[0].FocusUserID == A {
		return relationship_enum.RelationFans
	}
	return relationship_enum.RelationFous
}

func CalcUserPatchRelationship(A uint, BList []uint) (m map[uint]relationship_enum.Relation) {
	var userFocusList []models.UserFocusModel
	global.Db.Find(&userFocusList,
		"(user_id = ? OR focus_user_id in ?) or (focus_user_id = ? OR  user_id in  ?)",
		A, BList, A, BList)

	m = make(map[uint]relationship_enum.Relation)
	for _, B := range BList {
		m[B] = relationship_enum.RelationStranger

		var count int
		for _, model := range userFocusList {
			if model.FocusUserID == B {
				m[B] = relationship_enum.RelationFous
				count++
			}

			if model.UserID == B {
				m[B] = relationship_enum.RelationFans
				count++
			}
		}
		if count == 2 {
			m[B] = relationship_enum.RelationFriends
		}
	}
	return
}
