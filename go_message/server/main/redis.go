package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

// 以后统一在main.go将其初始化，在这不会自动调用
// idleTimeout time.Duration 时间类型
func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		// 最大空闲连接数
		MaxIdle: maxIdle,
		// 表示和数据库的最大连接数，0表示没有限制
		MaxActive: maxActive,
		// 最大空闲时间
		IdleTimeout: idleTimeout,
		// 初始化连接的代码，连接哪个 ip 的 redis
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}
