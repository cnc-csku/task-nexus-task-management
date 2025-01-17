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

func (m *mongoWorkspaceRepo) FindWorkspaceMemberByWorkspaceIDAndUserID(ctx context.Context, workspaceID bson.ObjectID, userID bson.ObjectID) (*models.WorkspaceMember, error) {
	f := NewWorkspaceFilter()
	f.WithWorkspaceID(workspaceID)
	f.WithMemberUserID(userID)

	var result struct {
		Members []models.WorkspaceMember `bson:"members"`
	}

	err := m.collection.FindOne(ctx, f).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	for _, member := range result.Members {
		if member.UserID == userID {
			return &member, nil
		}
	}

	return nil, nil
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

func (m *mongoWorkspaceRepo) CreateWorkspaceMember(ctx context.Context, in *repositories.CreateWorkspaceMemberRequest) error {
	f := NewWorkspaceFilter()
	f.WithWorkspaceID(in.WorkspaceID)

	update := bson.M{
		"$push": bson.M{
			"members": models.WorkspaceMember{
				UserID: in.UserID,
				Name:   in.Name,
				Role:   in.Role,
			},
		},
	}

	_, err := m.collection.UpdateOne(ctx, f, update)
	return err
}

func (m *mongoWorkspaceRepo) Create(ctx context.Context, workspace *repositories.CreateWorkspaceRequest) (*models.Workspace, error) {
	workspaceModel := models.Workspace{
		ID:   bson.NewObjectID(),
		Name: workspace.Name,
		Members: []models.WorkspaceMember{
			{
				UserID:    workspace.UserID,
				Name:      workspace.UserName,
				Role:      models.WorkspaceMemberRoleAdmin,
				JoinedAt:  time.Now(),
				RemovedAt: nil,
			},
		},
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
