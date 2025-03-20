package cache

import (
	"context"
	"log"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, config *config.Config) *redis.Client {
	opts, err := redis.ParseURL(config.Redis.URI)
	if err != nil {
		log.Fatalln("Redis URI is invalid: ", err)
	}

	client := redis.NewClient(opts)

	// test connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Fatalln("🚫 Cannot connect to Redis | ", err)
	} else {
		log.Println("✅ Connected to Redis")
	}

	return client
}
