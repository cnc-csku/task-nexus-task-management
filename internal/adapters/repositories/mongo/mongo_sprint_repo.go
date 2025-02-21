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

type mongoSprintRepo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoSprintRepo(config *config.Config, mongoClient *mongo.Client) repositories.SprintRepository {
	return &mongoSprintRepo{
		client:     mongoClient,
		collection: mongoClient.Database(config.MongoDB.Database).Collection("sprints"),
	}
}

func (m *mongoSprintRepo) Create(ctx context.Context, sprint *repositories.CreateSprintRequest) (*models.Sprint, error) {
	newSprint := models.Sprint{
		ID:        bson.NewObjectID(),
		ProjectID: sprint.ProjectID,
		Title:     sprint.Title,
		CreatedAt: time.Now(),
		CreatedBy: sprint.CreatedBy,
		UpdatedAt: time.Now(),
		UpdatedBy: sprint.CreatedBy,
	}

	_, err := m.collection.InsertOne(ctx, newSprint)
	if err != nil {
		return nil, err
	}

	return &newSprint, nil
}

func (m *mongoSprintRepo) FindByID(ctx context.Context, sprintID bson.ObjectID) (*models.Sprint, error) {
	sprint := new(models.Sprint)

	f := NewSprintFilter()
	f.WithID(sprintID)

	err := m.collection.FindOne(ctx, f).Decode(sprint)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return sprint, nil
}

func (m *mongoSprintRepo) Update(ctx context.Context, sprint *repositories.UpdateSprintRequest) error {
	f := NewSprintFilter()
	f.WithID(sprint.ID)

	u := bson.M{
		"$set": bson.M{
			"title":       sprint.Title,
			"sprint_goal": sprint.SprintGoal,
			"start_date":  sprint.StartDate,
			"end_date":    sprint.EndDate,
			"updated_at":  time.Now(),
			"updated_by":  sprint.UpdatedBy,
		},
	}

	_, err := m.collection.UpdateOne(ctx, f, u)
	if err != nil {
		return err
	}

	return nil
}
