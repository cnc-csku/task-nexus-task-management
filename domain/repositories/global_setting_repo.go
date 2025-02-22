package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
)

type GlobalSettingRepository interface {
	GetByKey(ctx context.Context, key string) (*models.KeyValuePair, error)
	Set(ctx context.Context, setting *models.KeyValuePair) error
}
