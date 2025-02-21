package cache

import (
	"log"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(config *config.Config) *redis.Client {
	opts, err := redis.ParseURL(config.Redis.URI)
	if err != nil {
		log.Fatalln("Redis URI is invalid: ", err)
	}

	return redis.NewClient(opts)
}
