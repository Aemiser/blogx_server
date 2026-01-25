package global

import (
	"blogx_server/conf"

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
)
