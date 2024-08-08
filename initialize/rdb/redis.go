package rdb

import "github.com/redis/go-redis/v9"

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "",
		Password: "", // 没有密码，默认值
		DB:       0,  // 默认DB 0
	})
}

func Rdb() *redis.Client {
	return rdb
}
