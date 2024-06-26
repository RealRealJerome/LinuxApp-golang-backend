package redisUtil

import (
	"github.com/go-redis/redis"
)

var RDB *redis.Client

// 初始化连接
func init() {
	// 通过 redis.NewClient 函数即可创建一个 redis 客户端, 这个方法接收一个 redis.Options 对象参数, 通过这个参数, 我们可以配置 redis 相关的属性, 例如 redis 服务器地址, 数据库名, 数据库密码等。
	RDB = redis.NewClient(&redis.Options{
		Addr:     "116.196.123.162:35787",
		Password: "R&ed!s@#R00Ot396&", // no password set
		DB:       0,                   // use default DB
	})
	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := RDB.Ping().Result()
	if err != nil {
		panic(err)
	}
}
