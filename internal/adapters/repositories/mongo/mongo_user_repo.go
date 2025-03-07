package mongo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"

	// "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
		ID:           bson.NewObjectID(),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		DisplayName:  user.DisplayName,
		ProfileUrl:   user.ProfileUrl,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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

func (m *mongoUserRepo) FindByIDs(ctx context.Context, userIDs []bson.ObjectID) ([]models.User, error) {
	f := NewUserFilter()
	f.WithUserIDs(userIDs)

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users := make([]models.User, 0)
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *mongoUserRepo) Search(ctx context.Context, in *repositories.SearchUserRequest) ([]*models.User, int64, error) {
	filter := bson.M{}
	if in.Keyword != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"email": bson.M{"$regex": in.Keyword, "$options": "i"}},
				{"full_name": bson.M{"$regex": in.Keyword, "$options": "i"}},
				{"display_name": bson.M{"$regex": in.Keyword, "$options": "i"}},
			},
		}
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((in.PaginationRequest.Page - 1) * in.PaginationRequest.PageSize))
	findOptions.SetLimit(int64(in.PaginationRequest.PageSize))

	// Set sorting options
	sortOrder := 1
	if strings.ToUpper(in.PaginationRequest.Order) == "DESC" {
		sortOrder = -1
	}
	findOptions.SetSort(bson.D{{Key: in.PaginationRequest.SortBy, Value: sortOrder}})

	cursor, err := m.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	users := make([]*models.User, 0)
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	// Get the total count of documents
	totalCount, err := m.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

func (m *mongoUserRepo) SearchWithUserIDs(ctx context.Context, in *repositories.SearchUserWithUserIDsRequest) ([]*models.User, int64, error) {
	filter := bson.M{
		"_id": bson.M{"$in": in.UserIDs},
	}

	if in.Keyword != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"email": bson.M{"$regex": in.Keyword, "$options": "i"}},
				{"full_name": bson.M{"$regex": in.Keyword, "$options": "i"}},
				{"display_name": bson.M{"$regex": in.Keyword, "$options": "i"}},
			},
		}
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((in.PaginationRequest.Page - 1) * in.PaginationRequest.PageSize))
	findOptions.SetLimit(int64(in.PaginationRequest.PageSize))

	// Set sorting options
	sortOrder := 1
	if strings.ToUpper(in.PaginationRequest.Order) == "DESC" {
		sortOrder = -1
	}
	findOptions.SetSort(bson.D{{Key: in.PaginationRequest.SortBy, Value: sortOrder}})

	cursor, err := m.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	users := make([]*models.User, 0)
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	// Get the total count of documents
	totalCount, err := m.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

func (m *mongoUserRepo) FindByID(ctx context.Context, userID bson.ObjectID) (*models.User, error) {
	user := new(models.User)

	f := NewUserFilter()
	f.WithUserID(userID)

	err := m.collection.FindOne(ctx, f).Decode(user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
