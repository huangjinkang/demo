package repositories

import (
	"context"
	"demo/src/common/redis_keys"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
	"time"
)

type RedisRepository struct {
	rdb *redis.Client
}

func NewRedisRepository(rdb *redis.Client) *RedisRepository {
	return &RedisRepository{rdb: rdb}
}

// LockArticleID 设置文章锁
func (repo *RedisRepository) LockArticleID(ctx context.Context, articleID uint64) (bool, error) {
	lockKey := redis_keys.GetArticleIdLockedKey(articleID)
	// 设置锁（这里设置锁的过期时间为3分钟，防止永久锁定）
	ok, err := repo.rdb.SetNX(ctx, lockKey, "locked", 3*time.Minute).Result()
	if err != nil {
		return false, err
	}
	return ok, nil // 如果 ok 为 true，则成功获取锁
}

// UnlockArticleID 释放文章锁
func (repo *RedisRepository) UnlockArticleID(ctx context.Context, articleID uint64) error {
	lockKey := redis_keys.GetArticleIdLockedKey(articleID)
	// 释放锁
	_, err := repo.rdb.Del(ctx, lockKey).Result()
	return err
}

// InitRedis 初始化Redis连接
func InitRedis() *redis.Client {
	// 从.env配置文件中获取Redis连接配置
	redisAddr := os.Getenv("REDIS_ADDR") // HOST:PORT
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")

	// 转换Redis数据库编号为int
	dbNum, err := strconv.Atoi(redisDB)
	if err != nil {
		log.Fatalf("Invalid Redis DB number: %v", err)
	}

	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       dbNum,
	})

	// 发送PING命令检查Redis是否连接成功
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Redis connected")
	return rdb
}
