package redis_site

import (
	"blogx_server/global"
	"context"
)

const (
	key = "blogx_site_flow"
)

func SetFlow() {
	v, _ := global.Redis.Get(context.Background(), key).Int()
	global.Redis.Set(context.Background(), key, v+1, 0)
}

func GetFlow() int {
	v, _ := global.Redis.Get(context.Background(), key).Int()
	return v
}

func Clean() {
	global.Redis.Del(context.Background(), key)
}
