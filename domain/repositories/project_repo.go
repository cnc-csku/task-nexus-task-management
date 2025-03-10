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
	UpdatePositions(ctx context.Context, projectID bson.ObjectID, position []string) error
	FindPositionByProjectID(ctx context.Context, projectID bson.ObjectID) ([]string, error)
	UpdateWorkflows(ctx context.Context, projectID bson.ObjectID, workflows []models.ProjectWorkflow) error
	FindWorkflowByProjectID(ctx context.Context, projectID bson.ObjectID) ([]models.ProjectWorkflow, error)
	IncrementSprintRunningNumber(ctx context.Context, projectID bson.ObjectID) error
	IncrementTaskRunningNumber(ctx context.Context, projectID bson.ObjectID) error
	UpdateAttributeTemplates(ctx context.Context, projectID bson.ObjectID, attributeTemplates []models.ProjectAttributeTemplate) error
	FindAttributeTemplatesByProjectID(ctx context.Context, projectID bson.ObjectID) ([]models.ProjectAttributeTemplate, error)
	UpdateSetupStatus(ctx context.Context, in *UpdateProjectSetupStatus) (*models.Project, error)
}

type CreateProjectRequest struct {
	WorkspaceID        bson.ObjectID
	Name               string
	ProjectPrefix      string
	Description        string
	Status             models.ProjectStatus
	Owner              *models.ProjectMember
	Workflows          []models.ProjectWorkflow
	AttributeTemplates []models.ProjectAttributeTemplate
	Positions          []string
	CreatedBy          bson.ObjectID
}

type UpdateProjectSetupStatus struct {
	ProjectID   bson.ObjectID
	SetupStatus models.ProjectSetupStatus
}
