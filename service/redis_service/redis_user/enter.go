package redis_user

import (
	"blogx_server/global"
	"context"
	"strconv"

	"github.com/sirupsen/logrus"
)

type userCacheType string

const (
	userCacheLook userCacheType = "user_look_key"
)

func set(t userCacheType, userID uint, n int) {
	num, _ := global.Redis.HGet(context.Background(), string(t), strconv.Itoa(int(userID))).Int()
	num += n
	global.Redis.HSet(context.Background(), string(t), strconv.Itoa(int(userID)), num)
}
func SetCacheLook(userID uint, increase bool) {
	var n = 1
	if !increase {
		n = -1
	}
	set(userCacheLook, userID, n)
}

func Clean() {
	global.Redis.Del(context.Background(), string(userCacheLook))
}
func get(t userCacheType, userID uint) int {
	num, _ := global.Redis.HGet(context.Background(), string(t), strconv.Itoa(int(userID))).Int()
	return num
}

func GetUserCacheLook(userID uint) int {
	return get(userCacheLook, userID)
}

func GetAll(artile userCacheType) (mps map[uint]int) {
	res, err := global.Redis.HGetAll(context.Background(), string(artile)).Result()
	if err != nil {
		return
	}
	mps = make(map[uint]int)
	for key, numS := range res {
		iK, err := strconv.Atoi(key)
		if err != nil {
			// 跳过无效的 key（可能是脏数据）
			logrus.Warnf("跳过无效的 article key: %s, value: %s", key, numS)
			continue
		}

		iN, err := strconv.Atoi(numS)
		if err != nil {
			// 跳过无效的值
			logrus.Warnf("跳过无效的 article value: key: %d, value: %s", iK, numS)
			continue
		}
		mps[uint(iK)] = iN
	}
	return mps
}

func GetAllCacheLook() (mps map[uint]int) {
	return GetAll(userCacheLook)
}
