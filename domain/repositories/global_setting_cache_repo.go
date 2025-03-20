package repositories

import (
	"context"
)

type KeyValuePairRedis struct {
	Key   string `redis:"key"`
	Type  string `redis:"type"`
	Value any    `redis:"value"`
}

type GlobalSettingCacheRepository interface {
	GetByKey(ctx context.Context, key string) (*KeyValuePairRedis, error)
	Set(ctx context.Context, setting *KeyValuePairRedis) error
}
