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

type mongoWorkspaceMemberRepo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoWorkspaceMemberRepo(config *config.Config, mongoClient *mongo.Client) repositories.WorkspaceMemberRepository {
	return &mongoWorkspaceMemberRepo{
		client:     mongoClient,
		collection: mongoClient.Database(config.MongoDB.Database).Collection("workspace_members"),
	}
}

func (m *mongoWorkspaceMemberRepo) FindByWorkspaceIDAndUserID(ctx context.Context, workspaceID bson.ObjectID, userID bson.ObjectID) (*models.WorkspaceMember, error) {
	f := NewWorkspaceMemberFilter()
	f.WithWorkspaceID(workspaceID)
	f.WithUserID(userID)

	var workspaceMember models.WorkspaceMember

	err := m.collection.FindOne(ctx, f).Decode(&workspaceMember)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &workspaceMember, nil
}

func (m *mongoWorkspaceMemberRepo) Create(ctx context.Context, in *repositories.CreateWorkspaceMemberRequest) error {
	workspaceMember := models.WorkspaceMember{
		ID:          bson.NewObjectID(),
		UserID:      in.UserID,
		WorkspaceID: in.WorkspaceID,
		Role:        in.Role,
		JoinedAt:    time.Now(),
	}

	_, err := m.collection.InsertOne(ctx, workspaceMember)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoWorkspaceMemberRepo) FindByUserID(ctx context.Context, userID bson.ObjectID) ([]models.WorkspaceMember, error) {
	f := NewWorkspaceMemberFilter()
	f.WithUserID(userID)

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	var workspaceMembers []models.WorkspaceMember
	if err = cursor.All(ctx, &workspaceMembers); err != nil {
		return nil, err
	}

	return workspaceMembers, nil
}

func (m *mongoWorkspaceMemberRepo) FindByWorkspaceID(ctx context.Context, workspaceID bson.ObjectID) ([]models.WorkspaceMember, error) {
	f := NewWorkspaceMemberFilter()
	f.WithWorkspaceID(workspaceID)

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	var workspaceMembers []models.WorkspaceMember
	if err = cursor.All(ctx, &workspaceMembers); err != nil {
		return nil, err
	}

	return workspaceMembers, nil
}
