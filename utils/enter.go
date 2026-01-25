package utils

func InList[T comparable](list []T, value T) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
