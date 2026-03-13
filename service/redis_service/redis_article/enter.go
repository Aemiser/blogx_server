package redis_article

import (
	"blogx_server/global"
	"context"
	"strconv"
)

type articleCacheType string

const (
	articleCacheDigg    articleCacheType = "article_digg_key"
	articleCacheCollect articleCacheType = "article_collect_key"
	articleCacheLook    articleCacheType = "article_look_key"
)

func set(t articleCacheType, articleID uint, increase bool) {
	num, _ := global.Redis.HGet(context.Background(), string(t), strconv.Itoa(int(articleID))).Int()
	if increase {
		num++
	} else {
		num--

	}
	global.Redis.HSet(context.Background(), string(t), strconv.Itoa(int(articleID)), num)
}
func SetCacheDigg(articleID uint, increase bool) {
	set(articleCacheDigg, articleID, increase)
}
func SetCacheCollect(articleID uint, increase bool) {
	set(articleCacheCollect, articleID, increase)
}
func SetCacheLook(articleID uint, increase bool) {
	set(articleCacheLook, articleID, increase)
}

func get(t articleCacheType, articleID uint) int {
	num, _ := global.Redis.HGet(context.Background(), string(t), strconv.Itoa(int(articleID))).Int()
	return num
}
func GetArticleCacheDigg(articleID uint) int {
	return get(articleCacheDigg, articleID)
}
func GetArticleCacheLook(articleID uint) int {
	return get(articleCacheLook, articleID)
}
func GetArticleCacheCollect(articleID uint) int {
	return get(articleCacheCollect, articleID)
}

func GetAll(artile articleCacheType) (mps map[uint]int) {
	res, err := global.Redis.HGetAll(context.Background(), string(artile)).Result()
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
func GetAllCacheDigg() (mps map[uint]int) {
	return GetAll(articleCacheDigg)
}
func GetAllCacheLook() (mps map[uint]int) {
	return GetAll(articleCacheLook)
}
func GetAllCacheCollect() (mps map[uint]int) {
	return GetAll(articleCacheCollect)
}
