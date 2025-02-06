package redis

import (
	"github.com/go-redis/redis_rate/v10"
	"github.com/k0lyaka/pow-antiddos/internal/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var Limiter *redis_rate.Limiter

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr:     config.Config.RedisHost,
		Password: "",
		DB:       0,
	})
	Limiter = redis_rate.NewLimiter(Client)
}
