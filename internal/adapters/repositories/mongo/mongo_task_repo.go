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

type mongoTaskRepo struct {
	config     *config.Config
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoTaskRepo(config *config.Config, mongoClient *mongo.Client) repositories.TaskRepository {
	return &mongoTaskRepo{
		config:     config,
		client:     mongoClient,
		collection: mongoClient.Database(config.MongoDB.Database).Collection("tasks"),
	}
}

func (m *mongoTaskRepo) Create(ctx context.Context, task *repositories.CreateTaskRequest) (*models.Task, error) {
	newTask := models.Task{
		ID:          bson.NewObjectID(),
		TaskID:      task.TaskID,
		ProjectID:   task.ProjectID,
		Title:       task.Title,
		Description: task.Description,
		ParentID:    task.ParentID,
		Type:        task.Type,
		Status:      task.Status,
		Sprint:      task.Sprint,
		CreatedAt:   time.Now(),
		CreatedBy:   task.CreatedBy,
		UpdatedAt:   time.Now(),
		UpdatedBy:   task.CreatedBy,
	}

	_, err := m.collection.InsertOne(ctx, newTask)
	if err != nil {
		return nil, err
	}

	return &newTask, nil
}

func (m *mongoTaskRepo) FindByID(ctx context.Context, id bson.ObjectID) (*models.Task, error) {
	task := new(models.Task)

	f := NewTaskFilter()
	f.WithID(id)

	err := m.collection.FindOne(ctx, f).Decode(task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return task, nil
}

func (m *mongoTaskRepo) FindByTaskID(ctx context.Context, taskID string) (*models.Task, error) {
	task := new(models.Task)

	f := NewTaskFilter()
	f.WithTaskID(taskID)

	err := m.collection.FindOne(ctx, f).Decode(task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return task, nil
}
