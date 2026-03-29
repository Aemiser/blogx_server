package cron_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/redis_service/redis_site"
)

// SyncSiteFlow 同步网站流量
func SyncSiteFlow() {
	flow := redis_site.GetFlow()
	global.Db.Create(&models.SiteFlowModel{Count: flow})
	redis_site.Clean()
}
