package cache

import (
	"os"
	"sync"

	"github.com/go-redis/redis"
)

type singleton struct {
	redis *redis.Client
}

var (
	once sync.Once

	instance singleton
)

func Redis() *redis.Client {
	once.Do(func() {
		instance = singleton{
			redis: redis.NewClient(&redis.Options{
				Addr:     os.Getenv("REDIS_URL"),
				Password: "", // no password set
				DB:       0,  // use default DB
			}),
		}
	})

	return instance.redis
}
