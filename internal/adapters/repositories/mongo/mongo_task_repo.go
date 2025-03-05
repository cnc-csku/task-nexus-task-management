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
		TaskRef:     task.TaskRef,
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

func (m *mongoTaskRepo) FindByTaskRef(ctx context.Context, taskRef string) (*models.Task, error) {
	task := new(models.Task)

	f := NewTaskFilter()
	f.WithTaskRef(taskRef)

	err := m.collection.FindOne(ctx, f).Decode(task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return task, nil
}

func (m *mongoTaskRepo) UpdateDetail(ctx context.Context, in *repositories.UpdateTaskDetailRequest) (*models.Task, error) {
	f := NewTaskFilter()
	f.WithID(in.ID)

	u := NewTaskUpdate()
	u.UpdateDetail(in)

	err := m.collection.FindOneAndUpdate(ctx, f, u).Err()
	if err != nil {
		return nil, err
	}

	return m.FindByID(ctx, in.ID)
}

func (m *mongoTaskRepo) UpdateStatus(ctx context.Context, in *repositories.UpdateTaskStatusRequest) (*models.Task, error) {
	f := NewTaskFilter()
	f.WithID(in.ID)

	u := NewTaskUpdate()
	u.UpdateStatus(in)

	err := m.collection.FindOneAndUpdate(ctx, f, u).Err()
	if err != nil {
		return nil, err
	}

	return m.FindByID(ctx, in.ID)
}

func (m *mongoTaskRepo) UpdateApprovals(ctx context.Context, in *repositories.UpdateTaskApprovalsRequest) (*models.Task, error) {
	f := NewTaskFilter()
	f.WithID(in.ID)

	u := NewTaskUpdate()
	u.UpdateApprovals(in)

	err := m.collection.FindOneAndUpdate(ctx, f, u).Err()
	if err != nil {
		return nil, err
	}

	return m.FindByID(ctx, in.ID)
}

func (m *mongoTaskRepo) ApproveTask(ctx context.Context, in *repositories.ApproveTaskRequest) (*models.Task, error) {
	f := NewTaskFilter()
	f.WithID(in.ID)
	f.WithUserApproval(in.UserID)

	u := NewTaskUpdate()
	u.ApproveTask(in.Reason)

	err := m.collection.FindOneAndUpdate(ctx, f, u).Err()
	if err != nil {
		return nil, err
	}

	return m.FindByID(ctx, in.ID)
}

func (m *mongoTaskRepo) UpdateAssignees(ctx context.Context, in *repositories.UpdateTaskAssigneesRequest) (*models.Task, error) {
	f := NewTaskFilter()
	f.WithID(in.ID)

	u := NewTaskUpdate()
	u.UpdateAssignees(in)

	err := m.collection.FindOneAndUpdate(ctx, f, u).Err()
	if err != nil {
		return nil, err
	}

	return m.FindByID(ctx, in.ID)
}

func (m *mongoTaskRepo) UpdateSprint(ctx context.Context, in *repositories.UpdateTaskSprintRequest) (*models.Task, error) {
	f := NewTaskFilter()
	f.WithID(in.ID)

	u := NewTaskUpdate()
	u.UpdateSprint(in)

	err := m.collection.FindOneAndUpdate(ctx, f, u).Err()
	if err != nil {
		return nil, err
	}

	return m.FindByID(ctx, in.ID)
}
