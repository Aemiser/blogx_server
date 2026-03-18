package redis_comment

import (
	"blogx_server/global"
	"context"
	"strconv"

	"github.com/sirupsen/logrus"
)

type commentCacheType string

const (
	commentCacheApply commentCacheType = "comment_apply_key"
	commentCacheDigg  commentCacheType = "comment_digg_key"
)

func set(t commentCacheType, commentID uint, n int) {
	num, _ := global.Redis.HGet(context.Background(), string(t), strconv.Itoa(int(commentID))).Int()
	num += n
	global.Redis.HSet(context.Background(), string(t), strconv.Itoa(int(commentID)), num)
}
func SetCacheApply(commentID uint, n int) {
	set(commentCacheApply, commentID, n)
}
func SetCacheDigg(commentID uint, n int) {
	set(commentCacheDigg, commentID, n)
}

func Clean() {
	global.Redis.Del(context.Background(), string(commentCacheApply))
}
func get(t commentCacheType, articleID uint) int {
	num, _ := global.Redis.HGet(context.Background(), string(t), strconv.Itoa(int(articleID))).Int()
	return num
}
func GetCacheApply(commentID uint) int {
	return get(commentCacheApply, commentID)
}
func GetCacheDigg(commentID uint) int {
	return get(commentCacheDigg, commentID)
}

func getAll(comment commentCacheType) (mps map[uint]int) {
	res, err := global.Redis.HGetAll(context.Background(), string(comment)).Result()
	if err != nil {
		return
	}
	mps = make(map[uint]int)
	for key, numS := range res {
		iK, err := strconv.Atoi(key)
		if err != nil {
			// 跳过无效的 key（可能是脏数据）
			logrus.Warnf("跳过无效的 comment key: %s, value: %s", key, numS)
			continue
		}

		iN, err := strconv.Atoi(numS)
		if err != nil {
			// 跳过无效的值
			logrus.Warnf("跳过无效的 comment value: key: %d, value: %s", iK, numS)
			continue
		}
		mps[uint(iK)] = iN
	}
	return mps
}
func GetAllCacheApply() (mps map[uint]int) {
	return getAll(commentCacheApply)
}
func GetAllCacheDigg() (mps map[uint]int) {
	return getAll(commentCacheDigg)
}
