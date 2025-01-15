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

type mongoProjectRepo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoProjectRepo(config *config.Config, mongoClient *mongo.Client) repositories.ProjectRepository {
	return &mongoProjectRepo{
		client:     mongoClient,
		collection: mongoClient.Database(config.MongoDB.Database).Collection("projects"),
	}
}

func (m *mongoProjectRepo) FindByWorkspaceIDAndName(ctx context.Context, workspaceID bson.ObjectID, name string) (*models.Project, error) {
	project := new(models.Project)

	f := NewProjectFilter()
	f.WithWorkspaceID(workspaceID)
	f.WithName(name)

	err := m.collection.FindOne(ctx, f).Decode(project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return project, nil
}

func (m *mongoProjectRepo) FindByWorkspaceIDAndProjectPrefix(ctx context.Context, workspaceID bson.ObjectID, projectPrefix string) (*models.Project, error) {
	project := new(models.Project)

	f := NewProjectFilter()
	f.WithWorkspaceID(workspaceID)
	f.WithProjectPrefix(projectPrefix)

	err := m.collection.FindOne(ctx, f).Decode(project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return project, nil
}

func (m *mongoProjectRepo) Create(ctx context.Context, project *repositories.CreateProjectRequest) (*models.Project, error) {
	newProject := models.Project{
		ID:            bson.NewObjectID(),
		WorkspaceID:   project.WorkspaceID,
		Name:          project.Name,
		ProjectPrefix: project.ProjectPrefix,
		Description:   project.Description,
		Status:        project.Status,
		Members:       project.Members,
		CreatedAt:     time.Now(),
		CreatedBy:     project.CreatedBy,
		UpdatedAt:     time.Now(),
		UpdatedBy:     project.CreatedBy,
	}

	_, err := m.collection.InsertOne(ctx, newProject)
	if err != nil {
		return nil, err
	}

	return &newProject, nil
}
