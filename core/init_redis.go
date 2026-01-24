package core

import (
	"blogx_server/global"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func InitRedis() *redis.Client {
	addr := fmt.Sprintf("%s:%d",
		global.Config.Redis.Host,
		global.Config.Redis.Port,
	)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,                         // Redis 服务器地址
		Password: global.Config.Redis.Password, // 如果没有密码，留空
		DB:       global.Config.Redis.DB,       // 使用默认数据库
	})

	// 测试连接
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logrus.Fatal("Redis数据库连接失败：%s", err)

	}
	logrus.Infof("Redis数据库连接成功!")
	return rdb
}
