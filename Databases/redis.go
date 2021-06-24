package Databases

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var RedisPool *redis.Pool

func RedisPollInit() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     20,
		MaxActive:   0,
		Wait:        true,
		IdleTimeout: time.Duration(1) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			redis.DialDatabase(0)
			return c, err
		},
	}
}

func RedisInit() {
	RedisPool = RedisPollInit()
}

func RedisClose() {
	_ = RedisPool.Get().Close()
}
