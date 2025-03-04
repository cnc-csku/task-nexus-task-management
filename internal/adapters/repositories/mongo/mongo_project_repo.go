package mongo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/constant"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoProjectRepo struct {
	config     *config.Config
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoProjectRepo(config *config.Config, mongoClient *mongo.Client) repositories.ProjectRepository {
	return &mongoProjectRepo{
		config:     config,
		client:     mongoClient,
		collection: mongoClient.Database(config.MongoDB.Database).Collection("projects"),
	}
}

func (m *mongoProjectRepo) FindByProjectID(ctx context.Context, projectID bson.ObjectID) (*models.Project, error) {
	project := new(models.Project)

	f := NewProjectFilter()
	f.WithID(projectID)

	err := m.collection.FindOne(ctx, f).Decode(project)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return project, nil
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
	// session, err := m.client.StartSession()
	// if err != nil {
	// 	return nil, err
	// }
	// defer session.EndSession(ctx)

	var newProject models.Project
	// _, err = session.WithTransaction(
	// 	ctx,
	// 	func(ctx context.Context) (interface{}, error) {
	newProject = models.Project{
		ID:                  bson.NewObjectID(),
		WorkspaceID:         project.WorkspaceID,
		Name:                project.Name,
		ProjectPrefix:       project.ProjectPrefix,
		Description:         project.Description,
		Status:              project.Status,
		SprintRunningNumber: 1,
		TaskRunningNumber:   1,
		Workflows:           project.Workflows,
		AttributeTemplates:  []models.AttributeTemplate{},
		Positions:           []string{},
		CreatedAt:           time.Now(),
		CreatedBy:           project.CreatedBy,
		UpdatedAt:           time.Now(),
		UpdatedBy:           project.CreatedBy,
	}

	_, err := m.collection.InsertOne(ctx, newProject)
	if err != nil {
		return nil, err
	}

	projectMemberCollection := m.client.Database(m.config.MongoDB.Database).Collection("project_members")

	newProjectMember := models.ProjectMember{
		ID:        bson.NewObjectID(),
		UserID:    project.Owner.UserID,
		ProjectID: newProject.ID,
		Role:      project.Owner.Role,
		JoinedAt:  time.Now(),
	}

	_, err = projectMemberCollection.InsertOne(ctx, newProjectMember)
	if err != nil {
		return nil, err
	}

	// return nil, nil
	// 	},
	// )
	// if err != nil {
	// 	return nil, err
	// }

	return &newProject, nil
}

func (m *mongoProjectRepo) FindByProjectIDsAndWorkspaceID(ctx context.Context, projectIDs []bson.ObjectID, workspaceID bson.ObjectID) ([]*models.Project, error) {
	f := NewProjectFilter()
	f.WithWorkspaceID(workspaceID)
	f.WithIDs(projectIDs)

	cursor, err := m.collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	projects := []*models.Project{}
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (m *mongoProjectRepo) FindMemberByProjectIDAndUserID(ctx context.Context, projectID bson.ObjectID, userID bson.ObjectID) (*models.ProjectMember, error) {
	f := NewProjectFilter()
	f.WithID(projectID)
	f.WithUserID(userID)

	var result struct {
		Members []models.ProjectMember `bson:"members"`
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

func (m *mongoProjectRepo) FindPositionByProjectID(ctx context.Context, projectID bson.ObjectID) ([]string, error) {
	f := NewProjectFilter()
	f.WithID(projectID)

	var result struct {
		Positions []string `bson:"positions"`
	}

	err := m.collection.FindOne(ctx, f).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return result.Positions, nil
}

func (m *mongoProjectRepo) UpdatePositions(ctx context.Context, projectID bson.ObjectID, position []string) error {
	f := NewProjectFilter()
	f.WithID(projectID)

	update := NewProjectUpdate()
	update.UpdatePositions(position)

	_, err := m.collection.UpdateOne(ctx, f, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoProjectRepo) SearchProjectMember(ctx context.Context, in *repositories.SearchProjectMemberRequest) ([]models.ProjectMember, int64, error) {
	f := NewProjectFilter()
	f.WithID(in.ProjectID)

	pipeline := []bson.M{
		{"$match": f},
		{"$unwind": "$members"},
	}

	if in.Keyword != "" {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"members.display_name": bson.M{"$regex": in.Keyword, "$options": "i"},
			},
		})
	}

	sortOrder := 1
	if strings.ToUpper(in.PaginationRequest.Order) == constant.DESC {
		sortOrder = -1
	}

	pipeline = append(pipeline, bson.M{
		"$sort": bson.M{
			"members." + in.PaginationRequest.SortBy: sortOrder,
		},
	})

	pipeline = append(pipeline, bson.M{
		"$skip": (in.PaginationRequest.Page - 1) * in.PaginationRequest.PageSize,
	})

	pipeline = append(pipeline, bson.M{
		"$limit": in.PaginationRequest.PageSize,
	})

	pipeline = append(pipeline, bson.M{
		"$group": bson.M{
			"_id":     "$_id",
			"members": bson.M{"$push": "$members"},
		},
	})

	cursor, err := m.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var result struct {
		Members []models.ProjectMember `bson:"members"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return nil, 0, err
		}
	}

	// Count the total number of members that match the filter
	countPipeline := []bson.M{
		{"$match": f},
		{"$unwind": "$members"},
	}

	if in.Keyword != "" {
		countPipeline = append(countPipeline, bson.M{
			"$match": bson.M{
				"members.display_name": bson.M{"$regex": in.Keyword, "$options": "i"},
			},
		})
	}

	countPipeline = append(countPipeline, bson.M{
		"$count": "count",
	})

	cursor, err = m.collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var countResult struct {
		Count int64 `bson:"count"`
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(&countResult); err != nil {
			return nil, 0, err
		}
	}

	return result.Members, countResult.Count, nil
}

func (m *mongoProjectRepo) UpdateWorkflows(ctx context.Context, projectID bson.ObjectID, workflows []models.Workflow) error {
	f := NewProjectFilter()
	f.WithID(projectID)

	update := NewProjectUpdate()

	bsonWorkflows := make([]bson.M, len(workflows))
	for i, w := range workflows {
		if w.PreviousStatuses == nil {
			w.PreviousStatuses = []string{}
		}
		bsonWorkflows[i] = bson.M{
			"previous_statuses": w.PreviousStatuses,
			"status":            w.Status,
			"is_default":        w.IsDefault,
		}
	}
	update.UpdateWorkflows(bsonWorkflows)

	_, err := m.collection.UpdateOne(ctx, f, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoProjectRepo) FindWorkflowByProjectID(ctx context.Context, projectID bson.ObjectID) ([]models.Workflow, error) {
	f := NewProjectFilter()
	f.WithID(projectID)

	var result struct {
		Workflows []models.Workflow `bson:"workflows"`
	}

	err := m.collection.FindOne(ctx, f).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return result.Workflows, nil
}

func (m *mongoProjectRepo) IncrementSprintRunningNumber(ctx context.Context, projectID bson.ObjectID) error {
	f := NewProjectFilter()
	f.WithID(projectID)

	u := NewProjectUpdate()
	u.IncrementSprintRunningNumber()

	_, err := m.collection.UpdateOne(ctx, f, u)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoProjectRepo) IncrementTaskRunningNumber(ctx context.Context, projectID bson.ObjectID) error {
	f := NewProjectFilter()
	f.WithID(projectID)

	u := NewProjectUpdate()
	u.IncrementTaskRunningNumber()

	_, err := m.collection.UpdateOne(ctx, f, u)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoProjectRepo) UpdateAttributeTemplates(ctx context.Context, projectID bson.ObjectID, attributeTemplates []models.AttributeTemplate) error {
	f := NewProjectFilter()
	f.WithID(projectID)

	update := NewProjectUpdate()

	bsonAttributeTemplates := make([]bson.M, len(attributeTemplates))
	for i, a := range attributeTemplates {
		bsonAttributeTemplates[i] = bson.M{
			"name": a.Name,
			"type": a.Type,
		}
	}
	update.UpdateAttributeTemplates(bsonAttributeTemplates)

	_, err := m.collection.UpdateOne(ctx, f, update)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoProjectRepo) FindAttributeTemplatesByProjectID(ctx context.Context, projectID bson.ObjectID) ([]models.AttributeTemplate, error) {
	f := NewProjectFilter()
	f.WithID(projectID)

	var result struct {
		AttributeTemplates []models.AttributeTemplate `bson:"attributes_templates"`
	}

	err := m.collection.FindOne(ctx, f).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return result.AttributeTemplates, nil
}
