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

type mongoProjectMemberRepo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoProjectMemberRepo(config *config.Config, mongoClient *mongo.Client) repositories.ProjectMemberRepository {
	return &mongoProjectMemberRepo{
		client:     mongoClient,
		collection: mongoClient.Database(config.MongoDB.Database).Collection("project_members"),
	}
}

func (m *mongoProjectMemberRepo) Create(ctx context.Context, in *repositories.CreateProjectMemberRequest) error {
	projectMember := models.ProjectMember{
		ID:        bson.NewObjectID(),
		UserID:    in.UserID,
		ProjectID: in.ProjectID,
		Role:      in.Role,
		Position:  in.Position,
		JoinedAt:  time.Now(),
	}

	_, err := m.collection.InsertOne(ctx, projectMember)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoProjectMemberRepo) CreateMany(ctx context.Context, projectMembers []repositories.CreateProjectMemberRequest) error {
	var projectMembersModel []models.ProjectMember
	for _, pm := range projectMembers {
		projectMember := models.ProjectMember{
			ID:        bson.NewObjectID(),
			UserID:    pm.UserID,
			ProjectID: pm.ProjectID,
			Role:      pm.Role,
			Position:  pm.Position,
			JoinedAt:  time.Now(),
		}
		projectMembersModel = append(projectMembersModel, projectMember)
	}

	_, err := m.collection.InsertMany(ctx, projectMembersModel)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoProjectMemberRepo) FindByUserID(ctx context.Context, userID bson.ObjectID) ([]*models.ProjectMember, error) {
	f := NewProjectMemberFilter()
	f.WithUserID(userID)

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	var projectMembers []*models.ProjectMember
	if err := cursor.All(ctx, &projectMembers); err != nil {
		return nil, err
	}

	return projectMembers, nil
}

func (m *mongoProjectMemberRepo) FindByProjectID(ctx context.Context, projectID bson.ObjectID) ([]*models.ProjectMember, error) {
	f := NewProjectMemberFilter()
	f.WithProjectID(projectID)

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	var projectMembers []*models.ProjectMember
	if err := cursor.All(ctx, &projectMembers); err != nil {
		return nil, err
	}

	return projectMembers, nil
}

func (m *mongoProjectMemberRepo) FindByProjectIDAndUserID(ctx context.Context, projectID bson.ObjectID, userID bson.ObjectID) (*models.ProjectMember, error) {
	f := NewProjectMemberFilter()
	f.WithUserID(userID)
	f.WithProjectID(projectID)

	projectMember := new(models.ProjectMember)
	err := m.collection.FindOne(ctx, f).Decode(projectMember)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return projectMember, nil
}
