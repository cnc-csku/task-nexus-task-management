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

type mongoWorkspaceRepo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoWorkspaceRepo(config *config.Config, mongoClient *mongo.Client) repositories.WorkspaceRepository {
	return &mongoWorkspaceRepo{
		client:     mongoClient,
		collection: mongoClient.Database(config.MongoDB.Database).Collection("workspaces"),
	}
}

func (m *mongoWorkspaceRepo) FindByID(ctx context.Context, workspaceID bson.ObjectID) (*models.Workspace, error) {
	f := NewWorkspaceFilter()
	f.WithWorkspaceID(workspaceID)

	var workspace models.Workspace
	err := m.collection.FindOne(ctx, f).Decode(&workspace)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &workspace, nil
}

func (m *mongoWorkspaceRepo) Create(ctx context.Context, workspace *repositories.CreateWorkspaceRequest) (*models.Workspace, error) {
	workspaceModel := models.Workspace{
		ID:        bson.NewObjectID(),
		Name:      workspace.Name,
		CreatedBy: workspace.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := m.collection.InsertOne(ctx, workspaceModel)
	if err != nil {
		return nil, err
	}

	workspaceModel.ID = res.InsertedID.(bson.ObjectID)
	return &workspaceModel, nil
}

func (m *mongoWorkspaceRepo) FindByWorkspaceIDs(ctx context.Context, workspaceIDs []bson.ObjectID) ([]models.Workspace, error) {
	f := NewWorkspaceFilter()
	f.WithWorkspaceIDs(workspaceIDs)

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	var workspaces []models.Workspace
	err = cursor.All(ctx, &workspaces)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}
