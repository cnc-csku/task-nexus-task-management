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
		Priority:    models.TaskPriorityMedium,
		Sprint:      task.Sprint,
		StartDate:   task.StartDate,
		DueDate:     task.DueDate,
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

func (m *mongoTaskRepo) FindByIDs(ctx context.Context, ids []bson.ObjectID) ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)

	f := NewTaskFilter()
	f.WithIDs(ids)

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
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

func (m *mongoTaskRepo) FindByTaskRefAndProjectID(ctx context.Context, taskRef string, projectID bson.ObjectID) (*models.Task, error) {
	task := new(models.Task)

	f := NewTaskFilter()
	f.WithTaskRef(taskRef)
	f.WithProjectID(projectID)

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

func (m *mongoTaskRepo) UpdateTitle(ctx context.Context, in *repositories.UpdateTaskTitleRequest) (*models.Task, error) {
	f := NewTaskFilter()
	f.WithID(in.ID)

	u := NewTaskUpdate()
	u.UpdateTitle(in)

	err := m.collection.FindOneAndUpdate(ctx, f, u).Err()
	if err != nil {
		return nil, err
	}

	return m.FindByID(ctx, in.ID)
}

func (m *mongoTaskRepo) UpdateParentID(ctx context.Context, in *repositories.UpdateTaskParentIDRequest) (*models.Task, error) {
	f := NewTaskFilter()
	f.WithID(in.ID)

	u := NewTaskUpdate()
	u.UpdateParentID(in)

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

func (m *mongoTaskRepo) FindByProjectIDAndStatuses(ctx context.Context, projectID bson.ObjectID, statuses []string) ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)

	f := NewTaskFilter()
	f.WithProjectID(projectID)
	f.WithStatuses(statuses)

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (m *mongoTaskRepo) UpdateHasChildren(ctx context.Context, in *repositories.UpdateTaskHasChildrenRequest) (*models.Task, error) {
	f := NewTaskFilter()
	f.WithID(in.ID)

	u := NewTaskUpdate()
	u.UpdateHasChildren(in.HasChildren)

	err := m.collection.FindOneAndUpdate(ctx, f, u).Err()
	if err != nil {
		return nil, err
	}

	return m.FindByID(ctx, in.ID)
}

func (m *mongoTaskRepo) FindByParentID(ctx context.Context, parentID bson.ObjectID) ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)

	f := NewTaskFilter()
	f.WithParentID(parentID)

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (m *mongoTaskRepo) UpdateChildrenPoint(ctx context.Context, in *repositories.UpdateTaskChildrenPointRequest) (*models.Task, error) {
	f := NewTaskFilter()
	f.WithID(in.ID)

	u := NewTaskUpdate()
	u.UpdateChildrenPoint(in.ChildrenPoint)

	err := m.collection.FindOneAndUpdate(ctx, f, u).Err()
	if err != nil {
		return nil, err
	}

	return m.FindByID(ctx, in.ID)
}

func (m *mongoTaskRepo) FindByProjectIDAndType(ctx context.Context, projectID bson.ObjectID, taskType models.TaskType) ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)

	f := NewTaskFilter()
	f.WithProjectID(projectID)
	f.WithType(taskType)

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (m *mongoTaskRepo) Search(ctx context.Context, in *repositories.SearchTaskRequest) ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)

	f := NewTaskFilter()
	f.WithProjectID(in.ProjectID)
	f.WithTypes(in.TaskTypes)

	if in.SprintID != nil {
		f.WithCurrentSprintID(*in.SprintID)
	}

	if in.IsTaskWithNoSprint {
		f.WithNoSprintID()
	}

	if in.EpicTaskID != nil {
		f.WithParentID(*in.EpicTaskID)
	}

	if in.IsTaskWithNoEpic {
		f.WithNoParentID()
	}

	if len(in.UserIDs) > 0 {
		f.WithUserIDs(in.UserIDs)
	}

	if len(in.Positions) > 0 {
		f.WithPositions(in.Positions)
	}

	if len(in.Statuses) > 0 {
		f.WithStatuses(in.Statuses)
	}

	if in.SearchKeyword != nil {
		f.WithSearchKeyword(*in.SearchKeyword)
	}

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (m *mongoTaskRepo) UpdateAttributes(ctx context.Context, in *repositories.UpdateTaskAttributesRequest) (*models.Task, error) {
	f := NewTaskFilter()
	f.WithID(in.ID)

	u := NewTaskUpdate()
	u.UpdateAttributes(in)

	err := m.collection.FindOneAndUpdate(ctx, f, u).Err()
	if err != nil {
		return nil, err
	}

	return m.FindByID(ctx, in.ID)
}

func (m *mongoTaskRepo) FindBySprintID(ctx context.Context, sprintID bson.ObjectID) ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)

	f := NewTaskFilter()
	f.WithCurrentSprintID(sprintID)

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
