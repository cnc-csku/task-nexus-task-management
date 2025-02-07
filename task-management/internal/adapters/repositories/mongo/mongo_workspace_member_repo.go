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

func (m *mongoWorkspaceMemberRepo) FindByWorkspaceID(ctx context.Context, workspaceID bson.ObjectID) ([]models.WorkspaceMember, error) {
	f := NewWorkspaceMemberFilter()
	f.WithWorkspaceID(workspaceID)

	workspaaceMembers := []models.WorkspaceMember{}
	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &workspaaceMembers)
	if err != nil {
		return nil, err
	}

	return workspaaceMembers, nil
}

func (m *mongoWorkspaceMemberRepo) Create(ctx context.Context, req *repositories.CreateWorkspaceMemberRequest) (*models.WorkspaceMember, error) {
	workspaceMemberModel := models.WorkspaceMember{
		ID:          bson.NewObjectID(),
		WorkspaceID: req.WorkspaceID,
		UserID:      req.UserID,
		Role:        req.Role,
		JoinedAt:    time.Now(),
	}

	_, err := m.collection.InsertOne(ctx, workspaceMemberModel)
	if err != nil {
		return nil, err
	}

	return &workspaceMemberModel, nil
}

func (m *mongoWorkspaceMemberRepo) FindByUserID(ctx context.Context, userID bson.ObjectID) ([]models.WorkspaceMember, error) {
	f := NewWorkspaceMemberFilter()
	f.WithUserID(userID)

	workspaaceMembers := []models.WorkspaceMember{}
	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &workspaaceMembers)
	if err != nil {
		return nil, err
	}

	return workspaaceMembers, nil
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
