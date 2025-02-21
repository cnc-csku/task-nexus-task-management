package mongo

import (
	"context"
	"strings"
	"time"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/constant"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoInvitationRepo struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoInvitationRepo(config *config.Config, mongoClient *mongo.Client) repositories.InvitationRepository {
	return &mongoInvitationRepo{
		client:     mongoClient,
		collection: mongoClient.Database(config.MongoDB.Database).Collection("invitations"),
	}
}

func (m *mongoInvitationRepo) FindByID(ctx context.Context, id bson.ObjectID) (*models.Invitation, error) {
	f := NewInvitationFilter()
	f.WithID(id)

	var invitation models.Invitation
	err := m.collection.FindOne(ctx, f).Decode(&invitation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &invitation, nil
}

func (m *mongoInvitationRepo) FindByWorkspaceIDAndInviteeUserID(ctx context.Context, workspaceID bson.ObjectID, inviteeUserID bson.ObjectID) (*models.Invitation, error) {
	f := NewInvitationFilter()
	f.WithWorkspaceID(workspaceID)
	f.WithInviteeUserID(inviteeUserID)
	f.WithNotExpired()
	f.WithNotResponded()

	var invitation models.Invitation
	err := m.collection.FindOne(ctx, f).Decode(&invitation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &invitation, nil
}

func (m *mongoInvitationRepo) Create(ctx context.Context, invitation *repositories.CreateInvitationRequest) error {
	newInvitation := models.Invitation{
		ID:            bson.NewObjectID(),
		WorkspaceID:   invitation.WorkspaceID,
		InviteeUserID: invitation.InviteeUserID,
		Role:          invitation.Role,
		Status:        invitation.Status,
		ExpiredAt:     invitation.ExpiredAt,
		CustomMessage: &invitation.CustomMessage,
		CreatedAt:     time.Now(),
		CreatedBy:     invitation.CreatedBy,
	}

	_, err := m.collection.InsertOne(ctx, newInvitation)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoInvitationRepo) FindByInviteeUserID(ctx context.Context, inviteeUserID bson.ObjectID, sortBy string, order string) ([]models.Invitation, error) {
	f := NewInvitationFilter()
	f.WithInviteeUserID(inviteeUserID)

	findOptions := options.Find()
	sortOrder := 1
	if strings.ToUpper(order) == constant.DESC {
		sortOrder = -1
	}
	findOptions.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})

	cursor, err := m.collection.Find(ctx, f, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invitations []models.Invitation
	if err := cursor.All(ctx, &invitations); err != nil {
		return nil, err
	}

	return invitations, nil
}

func (m *mongoInvitationRepo) UpdateStatus(ctx context.Context, id bson.ObjectID, status models.InvitationStatus) error {
	f := NewInvitationFilter()
	f.WithID(id)

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: status},
			{Key: "responded_at", Value: time.Now()},
		}},
	}

	_, err := m.collection.UpdateOne(ctx, f, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoInvitationRepo) SearchInvitationForEachWorkspace(ctx context.Context, in *repositories.SearchInvitationForEachWorkspaceRequest) ([]models.Invitation, int64, error) {
	filter := bson.M{"workspace_id": in.WorkspaceID}

	if in.Keyword != "" {
		if in.SearchBy != "" {
			filter[in.SearchBy] = bson.M{"$regex": in.Keyword, "$options": "i"}
		} else {
			filter["$or"] = []bson.M{
				{"status": bson.M{"$regex": in.Keyword, "$options": "i"}},
				{"custom_message": bson.M{"$regex": in.Keyword, "$options": "i"}},
			}
		}
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((in.PaginationRequest.Page - 1) * in.PaginationRequest.PageSize))
	findOptions.SetLimit(int64(in.PaginationRequest.PageSize))

	sortOrder := 1
	if strings.ToUpper(in.PaginationRequest.Order) == constant.DESC {
		sortOrder = -1
	}
	findOptions.SetSort(bson.D{{Key: in.PaginationRequest.SortBy, Value: sortOrder}})

	cursor, err := m.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var invitations []models.Invitation
	for cursor.Next(ctx) {
		var invitation models.Invitation
		if err := cursor.Decode(&invitation); err != nil {
			return nil, 0, err
		}
		invitations = append(invitations, invitation)
	}

	total, err := m.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return invitations, total, nil
}
