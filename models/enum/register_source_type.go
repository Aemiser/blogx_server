package enum

type RegisterSourceType int8

const (
	RegisterSourceTypeEmail RegisterSourceType = iota + 1
	RegisterSourceTypePhone
	RegisterSourceTypeQQ
	RegisterSourceTypeWeChat
	RegisterSourceTypeCmd
)
