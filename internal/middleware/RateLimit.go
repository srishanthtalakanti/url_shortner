package middleware

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type RateLimit struct {
	RedisClient *redis.Client
	Limit       int
	Duration    time.Duration
}
