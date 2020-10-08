package redispool

import (
	"time"
	"github.com/gomodule/redigo/redis"
)

// NewPool f
func NewPool(server string) *redis.Pool {
	return &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,

			Dial: func() (redis.Conn, error) {
					c, err := redis.Dial("tcp", "redis:6379")

					if err != nil {
							return nil, err
					}

					return c, err
			},

			TestOnBorrow: func(c redis.Conn, t time.Time) error {
					_, err := c.Do("PING")
					return err
			},
	}
}