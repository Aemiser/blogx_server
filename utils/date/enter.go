package date

import "time"

func GetNowAfter() time.Time {
	now := time.Now()
	endTimeObj := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	return endTimeObj
}
