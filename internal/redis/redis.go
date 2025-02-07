package redis

import (
	"log"

	"github.com/go-redis/redis_rate/v10"
	"github.com/k0lyaka/pow-antiddos/internal/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var Limiter *redis_rate.Limiter

func InitRedis() {
	opt, err := redis.ParseURL(config.Config.RedisUrl)

	if err != nil {
		log.Fatalln(err)
	}

	Client = redis.NewClient(opt)
	Limiter = redis_rate.NewLimiter(Client)
}
