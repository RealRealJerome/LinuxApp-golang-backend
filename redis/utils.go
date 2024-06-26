package redisUtil

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

// SetKV 设置一个 KV 对并设置超时时间
func SetKV(key string, value interface{}, expiration time.Duration) error {
	err := RDB.Set(key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetKV 查询一个 KV 对
func GetKV(key string) (string, error) {
	val, err := RDB.Get(key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist")
	} else if err != nil {
		return "", err
	}
	return val, nil
}
