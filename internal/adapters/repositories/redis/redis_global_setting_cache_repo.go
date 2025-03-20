package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/redis/go-redis/v9"
)

const GLOBAL_KEY_FORMAT = "global-setting:%s"

type redisGlobalSettingCacheRepo struct {
	config *config.Config
	client *redis.Client
}

type KeyValuePairRedisString struct {
	Key   string `redis:"key"`
	Type  string `redis:"type"`
	Value string `redis:"value"`
}

func NewRedisGlobalSettingCacheRepo(config *config.Config, client *redis.Client) repositories.GlobalSettingCacheRepository {
	return &redisGlobalSettingCacheRepo{
		config: config,
		client: client,
	}
}

func (r *redisGlobalSettingCacheRepo) GetByKey(ctx context.Context, key string) (*repositories.KeyValuePairRedis, error) {

	var tmpResult KeyValuePairRedisString

	err := r.client.HGetAll(ctx, fmt.Sprintf(GLOBAL_KEY_FORMAT, key)).Scan(&tmpResult)
	if err != nil {
		return nil, err
	}

	result := repositories.KeyValuePairRedis{
		Key:   tmpResult.Key,
		Type:  tmpResult.Type,
		Value: tmpResult.Value,
	}

	// Check if the result is empty (key not found)
	if result.Key == "" {
		return nil, fmt.Errorf("global setting not found for key: %s", key)
	}

	return &result, nil
}

func (r *redisGlobalSettingCacheRepo) Set(ctx context.Context, setting *repositories.KeyValuePairRedis) error {
	ttl, err := time.ParseDuration(r.config.Cache.GlobalConfigTTL)
	if err != nil {
		return err
	}

	var val string

	switch setting.Type {
	case string(models.KeyValuePairTypeString):
		val = setting.Value.(string)
	case string(models.KeyValuePairTypeBoolean):
		val = fmt.Sprintf("%t", setting.Value.(bool))
	case string(models.KeyValuePairTypeNumber):
		val = fmt.Sprintf("%f", setting.Value.(float64))
	case string(models.KeyValuePairTypeDate):
		val = setting.Value.(time.Time).Format(time.RFC3339)
	}

	stringPair := KeyValuePairRedisString{
		Key:   setting.Key,
		Type:  setting.Type,
		Value: val,
	}

	err = r.client.HSet(ctx, fmt.Sprintf(GLOBAL_KEY_FORMAT, setting.Key), stringPair).Err()
	if err != nil {
		return err
	}

	err = r.client.Expire(ctx, fmt.Sprintf(GLOBAL_KEY_FORMAT, setting.Key), ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
