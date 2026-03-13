package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

func main() {
	now := time.Now()
	endTime := now.Format("2006-01-02") + " 23:59:59"
	endTimeObj, err := time.Parse("2006-01-02 15:04:05", endTime)

	if err != nil {
		logrus.Error("time parse err:", err)
	}
	subTime := endTimeObj.Sub(now)
	fmt.Println(endTimeObj)
	fmt.Println(subTime)
	fmt.Println(now)

	time2 := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	fmt.Println(time2)
	fmt.Println(time2.Sub(now))

}
