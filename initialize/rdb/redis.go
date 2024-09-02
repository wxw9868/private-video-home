package rdb

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/wxw9868/video/config"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {
	conf := config.Config().Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password, // 没有密码，默认值
		DB:       conf.DB,       // 默认DB 0
	})

	// 测试连接
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis连接超时: %s\n", err)
	}
	log.Printf("Redis连接成功: %s\n", pong)
}

func Rdb() *redis.Client {
	return rdb
}
