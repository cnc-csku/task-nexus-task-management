package repositories

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProjectRepository interface {
	FindByProjectID(ctx context.Context, projectID bson.ObjectID) (*models.Project, error)
	FindByWorkspaceIDAndName(ctx context.Context, workspaceID bson.ObjectID, name string) (*models.Project, error)
	FindByWorkspaceIDAndProjectPrefix(ctx context.Context, workspaceID bson.ObjectID, projectPrefix string) (*models.Project, error)
	Create(ctx context.Context, project *CreateProjectRequest) (*models.Project, error)
	FindByProjectIDsAndWorkspaceID(ctx context.Context, projectIDs []bson.ObjectID, workspaceID bson.ObjectID) ([]*models.Project, error)
	AddPositions(ctx context.Context, projectID bson.ObjectID, position []string) error
	FindPositionByProjectID(ctx context.Context, projectID bson.ObjectID) ([]string, error)
	AddWorkflows(ctx context.Context, projectID bson.ObjectID, workflows []models.Workflow) error
	FindWorkflowByProjectID(ctx context.Context, projectID bson.ObjectID) ([]models.Workflow, error)
	IncrementSprintRunningNumber(ctx context.Context, projectID bson.ObjectID) error
	IncrementTaskRunningNumber(ctx context.Context, projectID bson.ObjectID) error
	AddAttributeTemplates(ctx context.Context, projectID bson.ObjectID, attributeTemplates []models.AttributeTemplate) error
	FindAttributeTemplatesByProjectID(ctx context.Context, projectID bson.ObjectID) ([]models.AttributeTemplate, error)
}

type CreateProjectRequest struct {
	WorkspaceID   bson.ObjectID
	Name          string
	ProjectPrefix string
	Description   *string
	Status        models.ProjectStatus
	Owner         *models.ProjectMember
	Workflows     []models.Workflow
	CreatedBy     bson.ObjectID
}
