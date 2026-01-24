package enum

type RoleType int8

const (
	UserRole RoleType = iota + 1
	AdminRole
	visitorRole
)
