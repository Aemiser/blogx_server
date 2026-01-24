package redis_jwt

import (
	"blogx_server/common/jwts"
	"blogx_server/global"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RedisJwtKey string

const (
	RedisTokenBlack RedisJwtKey = "token:black:%s"
)

type BlackType int8

const (
	UserBlackType BlackType = iota + 1
	AdminBlackType
	DeviceBlackType
)

func (b BlackType) Msg() string {
	switch b {
	case UserBlackType:
		return "已注销"
	case AdminBlackType:
		return "禁止登入"
	case DeviceBlackType:
		return "设备下线"
	default:
		return "已注销"
	}
}
func (b BlackType) String() string {
	return fmt.Sprintf("%d", b)
}

func (b BlackType) ParseBlackType(value string) BlackType {
	switch value {
	case UserBlackType.String():
		return UserBlackType
	case AdminBlackType.String():
		return AdminBlackType
	case DeviceBlackType.String():
		return DeviceBlackType
	default:
		return UserBlackType
	}
}

func BlackToken(token string, value BlackType) {
	key := fmt.Sprintf(string(RedisTokenBlack), token)

	claim, err := jwts.ParseToken(token)
	if err != nil || claim == nil {
		logrus.Errorf("token解析失败:%s", err)
		return
	}
	second := claim.ExpiresAt - time.Now().Unix()
	res := global.Redis.Set(context.Background(), key, value.String(), time.Duration(second)*time.Second)
	if res.Err() != nil {
		logrus.Errorf("token加入黑名单失败:%s", res.Err())
	}
}

func HasTokenBlack(token string) (blk BlackType, ok bool) {
	key := fmt.Sprintf(string(RedisTokenBlack), token)
	val, err := global.Redis.Get(context.Background(), key).Result()
	if err != nil && val != "" {
		logrus.Errorf("token查询失败:val:%s,err:%s", val, err)
		return
	}

	if val == "" {
		return
	}
	blk = blk.ParseBlackType(val)
	return blk, true

}

func HasTokenBlackByGin(c *gin.Context) (blk BlackType, ok bool) {
	tokenString := c.GetHeader("token")
	if tokenString == "" {
		tokenString = c.Query("token")
	}
	return HasTokenBlack(tokenString)
}
