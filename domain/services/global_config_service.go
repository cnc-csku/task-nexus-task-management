package services

import (
	"context"
	"strconv"
	"time"

	"github.com/cnc-csku/task-nexus-go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type GlobalSettingService interface {
	GetGlobalSettingByKey(ctx context.Context, key string) (*models.KeyValuePair, *errutils.Error)
	SetGlobalSetting(ctx context.Context, setting *models.KeyValuePair) *errutils.Error
}

type globalSettingServiceImpl struct {
	globalSettingRepo      repositories.GlobalSettingRepository
	globalSettingCacheRepo repositories.GlobalSettingCacheRepository
}

func NewGlobalSettingService(
	globalSettingRepo repositories.GlobalSettingRepository,
	globalSettingCacheRepo repositories.GlobalSettingCacheRepository,
) GlobalSettingService {
	return &globalSettingServiceImpl{
		globalSettingRepo:      globalSettingRepo,
		globalSettingCacheRepo: globalSettingCacheRepo,
	}
}

func (g *globalSettingServiceImpl) GetGlobalSettingByKey(ctx context.Context, key string) (*models.KeyValuePair, *errutils.Error) {
	// Try to get from cache first
	cachedSetting, err := g.globalSettingCacheRepo.GetByKey(ctx, key)

	if err != nil {
		// If not in cache, fetch from repository
		setting, err := g.globalSettingRepo.GetByKey(ctx, key)
		if err != nil {
			return nil, errutils.NewError(err, errutils.InternalServerError)
		}

		if setting.Type == models.KeyValuePairTypeDate {
			setting.Value = setting.Value.(bson.DateTime).Time()
		}

		// Update cache
		err = g.globalSettingCacheRepo.Set(ctx, &repositories.KeyValuePairRedis{
			Key:   setting.Key,
			Type:  string(setting.Type),
			Value: setting.Value,
		})

		if err != nil {
			return nil, errutils.NewError(err, errutils.InternalServerError)
		}

		return setting, nil
	}

	var value any

	// Convert type
	switch cachedSetting.Type {
	case string(models.KeyValuePairTypeBoolean):
		value = cachedSetting.Value.(string) == "true"
	case string(models.KeyValuePairTypeString):
		value = cachedSetting.Value.(string)
	case string(models.KeyValuePairTypeNumber):
		number, err := strconv.ParseFloat(cachedSetting.Value.(string), 64)
		if err != nil {
			return nil, errutils.NewError(err, errutils.InternalServerError)
		}
		value = number
	case string(models.KeyValuePairTypeDate):
		t, err := time.Parse(time.RFC3339, cachedSetting.Value.(string))
		if err != nil {
			return nil, errutils.NewError(err, errutils.InternalServerError)
		}
		value = t
	}

	return &models.KeyValuePair{
		Key:   cachedSetting.Key,
		Type:  models.KeyValuePairType(cachedSetting.Type),
		Value: value,
	}, nil
}

func (g *globalSettingServiceImpl) SetGlobalSetting(ctx context.Context, setting *models.KeyValuePair) *errutils.Error {
	// Set to database
	err := g.globalSettingRepo.Set(ctx, setting)
	if err != nil {
		return errutils.NewError(err, errutils.InternalServerError)
	}

	// Set to cache
	// Ignore error it's can be set later
	_ = g.globalSettingCacheRepo.Set(ctx, &repositories.KeyValuePairRedis{
		Key:   setting.Key,
		Type:  string(setting.Type),
		Value: setting.Value,
	})

	return nil
}
