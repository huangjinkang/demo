package redis_keys

import "fmt"

const Lock = "lock"

// GetArticleIdLockedKey 获取文章锁的键名
func GetArticleIdLockedKey(articleId uint64) (lockedKey string) {
	lockedKey = fmt.Sprintf("%s:%d", Lock, articleId)
	return
}
