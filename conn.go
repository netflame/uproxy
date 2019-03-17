package main

import (
	"github.com/gomodule/redigo/redis"
)

var (
	redisPool *redis.Pool
)

// try to use this
func initRedisPool() {
	redisPool = getRedisPoolFromConfig(config)
}

func getConn() redis.Conn {
	return redisPool.Get()
}

func getRedisPool() *redis.Pool {
	c := getConfig()
	rp := getRedisPoolFromConfig(c)
	return rp
}

func getRedisPoolFromConfig(c *Config) *redis.Pool {
	// pc: redis Pool Config
	// rc: Redis Config
	// rp: Redis Pool
	pc, rc := c.RPool, c.Redis
	rp := &redis.Pool{
		MaxIdle:         pc.MaxIdle,
		MaxActive:       pc.MaxActive,
		IdleTimeout:     pc.IdleTimeout.Duration,
		Wait:            pc.Wait,
		MaxConnLifetime: pc.MaxConnLifetime.Duration,
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial(
				"tcp", rc.Address,
				redis.DialPassword(rc.Password),
				redis.DialDatabase(rc.DB),
				redis.DialConnectTimeout(rc.ConnectTimeout.Duration),
				redis.DialReadTimeout(rc.ReadTimeout.Duration),
				redis.DialWriteTimeout(rc.WriteTimeout.Duration),
			)
			if err != nil {
				return nil, err
			}
			return
		},
	}
	return rp
}
