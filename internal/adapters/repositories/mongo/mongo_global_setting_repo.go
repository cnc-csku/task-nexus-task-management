package mongo

import (
	"context"
	"errors"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoGlobalSettingRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoGlobalSettingRepo(config *config.Config, client *mongo.Client) repositories.GlobalSettingRepository {
	return &mongoGlobalSettingRepository{
		client:     client,
		collection: client.Database(config.MongoDB.Database).Collection("global_settings"),
	}
}

func (m *mongoGlobalSettingRepository) GetByKey(ctx context.Context, key string) (*models.KeyValuePair, error) {
	f := NewGlobalSettingFilter()
	f.WithKey(key)

	var setting models.KeyValuePair
	err := m.collection.FindOne(ctx, f).Decode(&setting)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &setting, nil
}

func (m *mongoGlobalSettingRepository) Set(ctx context.Context, setting *models.KeyValuePair) error {
	f := NewGlobalSettingFilter()
	f.WithKey(setting.Key)

	u := NewGlobalSettingUpdate()
	u.WithSetting(setting)

	_, err := m.collection.UpdateOne(ctx, f, u, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}
