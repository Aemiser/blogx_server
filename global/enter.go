package global

import (
	"blogx_server/conf"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	Version = "10.0.1"
)

var (
	Config *conf.Config
	Db     *gorm.DB
	Redis  *redis.Client
	// 库中提供的默认内存存储器。内存存储器用于存储和验证生成的验证码信息。它将验证码的标识符、验证码图片和相关的验证数据存储在内存中。
	Stores = base64Captcha.DefaultMemStore
)
