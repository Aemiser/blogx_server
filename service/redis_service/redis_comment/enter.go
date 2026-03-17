package redis_comment

import (
	"blogx_server/global"
	"context"
	"strconv"
)

type commentCacheType string

const (
	commentCacheApply commentCacheType = "comment_apply_key"
)

func set(t commentCacheType, commentID uint, n int) {
	num, _ := global.Redis.HGet(context.Background(), string(t), strconv.Itoa(int(commentID))).Int()
	num += n
	global.Redis.HSet(context.Background(), string(t), strconv.Itoa(int(commentID)), num)
}
func SetCacheApply(articleID uint, n int) {
	set(commentCacheApply, articleID, n)
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

func GetAll(comment commentCacheType) (mps map[uint]int) {
	res, err := global.Redis.HGetAll(context.Background(), string(comment)).Result()
	if err != nil {
		return
	}
	mps = make(map[uint]int)
	for key, numS := range res {
		iK, err := strconv.Atoi(key)
		if err != nil {
			continue
		}

		iN, err := strconv.Atoi(numS)
		if err != nil {
			continue
		}
		mps[uint(iK)] = iN
	}
	return mps
}
func GetAllCacheApply() (mps map[uint]int) {
	return GetAll(commentCacheApply)
}
