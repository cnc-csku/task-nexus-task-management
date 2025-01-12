package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoUserRepo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoUserRepo(config *config.Config, mongoClient *mongo.Client) repositories.UserRepository {
	return &mongoUserRepo{
		client:     mongoClient,
		collection: mongoClient.Database(config.MongoDB.Database).Collection("users"),
	}
}

func (m *mongoUserRepo) Create(ctx context.Context, user *repositories.CreateUserRequest) (*models.User, error) {
	newUser := models.User{
		ID: bson.NewObjectID(),
		Email: user.Email,
		PasswordHash: user.PasswordHash,
		FullName: user.FullName,
		DisplayName: user.DisplayName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := m.collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (m *mongoUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)

	f := NewUserFilter()
	f.WithEmail(email)

	err := m.collection.FindOne(ctx, f).Decode(user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

