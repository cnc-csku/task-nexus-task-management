package mongo

import (
	"context"
	"time"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoTaskCommentRepo struct {
	config     *config.Config
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoTaskCommentRepo(config *config.Config, mongoClient *mongo.Client) repositories.TaskCommentRepository {
	return &mongoTaskCommentRepo{
		config:     config,
		client:     mongoClient,
		collection: mongoClient.Database(config.MongoDB.Database).Collection("task_comments"),
	}
}

func (m *mongoTaskCommentRepo) Create(ctx context.Context, taskComment *repositories.CreateTaskCommentRequest) (*models.TaskComment, error) {
	newTaskComment := models.TaskComment{
		ID:        bson.NewObjectID(),
		Content:   taskComment.Content,
		UserID:    taskComment.UserID,
		TaskID:    taskComment.TaskID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := m.collection.InsertOne(ctx, newTaskComment)
	if err != nil {
		return nil, err
	}

	return &newTaskComment, nil
}

func (m *mongoTaskCommentRepo) FindByTaskID(ctx context.Context, taskID bson.ObjectID) ([]*models.TaskComment, error) {
	f := NewTaskCommentFilter()
	f.WithTaskID(taskID)

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	taskComments := make([]*models.TaskComment, 0)
	if err := cursor.All(ctx, &taskComments); err != nil {
		return nil, err
	}

	return taskComments, nil
}
