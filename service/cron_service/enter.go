package cron_service

import (
	"time"

	"github.com/robfig/cron/v3"
)

func Cron() {
	//crontab := cron.New()  默认从分开始进行时间调度
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))

	// 凌晨两点同步
	crontab.AddFunc("* * 2 * * *", SyncArticle)
	crontab.AddFunc("* * 3 * * *", SyncComment)
	crontab.AddFunc("50 59 23 * * *", SyncSiteFlow)
	crontab.Start()
}
