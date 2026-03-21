package relationship_enum

type Relation int8

const (
	RelationStranger Relation = iota + 1
	RelationFous
	RelationFans
	RelationFriends
)
