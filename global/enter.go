package global

import (
	"blogx_server/conf"
	"sync"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	Version = "10.0.1"
)

var (
	Config           *conf.Config
	Db               *gorm.DB
	Redis            *redis.Client
	Stores           = base64Captcha.DefaultMemStore
	EmailVerifyStore = sync.Map{}
)
