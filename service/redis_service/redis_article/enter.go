package redis_article

import (
	"blogx_server/global"
	"blogx_server/utils/date"
	"context"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

type articleCacheType string

const (
	articleCacheDigg    articleCacheType = "article_digg_key"
	articleCacheCollect articleCacheType = "article_collect_key"
	articleCacheLook    articleCacheType = "article_look_key"
	articleCacheComment articleCacheType = "article_comment_key"
)

func set(t articleCacheType, articleID uint, n int) {
	num, _ := global.Redis.HGet(context.Background(), string(t), strconv.Itoa(int(articleID))).Int()
	num += n
	global.Redis.HSet(context.Background(), string(t), strconv.Itoa(int(articleID)), num)
}
func SetCacheDigg(articleID uint, increase bool) {
	var n = 1
	if !increase {
		n = -1
	}
	set(articleCacheDigg, articleID, n)
}
func SetCacheCollect(articleID uint, increase bool) {
	var n = 1
	if !increase {
		n = -1
	}
	set(articleCacheCollect, articleID, n)
}
func SetCacheLook(articleID uint, increase bool) {
	var n = 1
	if !increase {
		n = -1
	}
	set(articleCacheLook, articleID, n)
}

func SetCacheComment(articleID uint, n int) {
	set(articleCacheComment, articleID, n)
}

func Clean() {
	global.Redis.Del(context.Background(), string(articleCacheDigg), string(articleCacheCollect), string(articleCacheLook))
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
func GetArticleCacheComment(articleID uint) int {
	return get(articleCacheComment, articleID)
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

func GetAllCacheComment() (mps map[uint]int) {
	return GetAll(articleCacheComment)
}

func SetUserArticleHistoryCache(articleID uint, userID uint) {
	key := fmt.Sprintf("history_%d", userID)
	field := fmt.Sprintf("%d", articleID)
	endTimeObj := date.GetNowAfter()
	err := global.Redis.HSet(context.Background(), key, field, "").Err()
	if err != nil {
		logrus.Error("redis set err:", err)
		return
	}

	err = global.Redis.ExpireAt(context.Background(), key, endTimeObj).Err()
	if err != nil {
		logrus.Error("redis expire err:", err)
		return
	}
}

func GetUserArticleHistoryCache(articleID uint, userID uint) bool {
	key := fmt.Sprintf("history_%d", userID)
	field := fmt.Sprintf("%d", articleID)
	err := global.Redis.HGet(context.Background(), key, field).Err()
	if err != nil {
		return false
	}
	return true
}
