package redis_jwt

import (
	"blogx_server/common/jwts"
	"blogx_server/global"
	"context"
	"fmt"
	"time"

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
	if err != nil {
		logrus.Errorf("token查询失败:%s", err)
		return
	}
	blk = blk.ParseBlackType(val)
	return blk, true

}
