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
	FindByWorkspaceIDAndUserID(ctx context.Context, workspaceID bson.ObjectID, userID bson.ObjectID) ([]*models.Project, error)
	FindMemberByProjectIDAndUserID(ctx context.Context, projectID bson.ObjectID, userID bson.ObjectID) (*models.ProjectMember, error)
	AddPositions(ctx context.Context, projectID bson.ObjectID, position []string) error
	FindPositionByProjectID(ctx context.Context, projectID bson.ObjectID) ([]string, error)
	AddMembers(ctx context.Context, projectID bson.ObjectID, member []CreateProjectMemberRequest) error
	SearchProjectMember(ctx context.Context, in *SearchProjectMemberRequest) ([]models.ProjectMember, int64, error)
	AddWorkflows(ctx context.Context, projectID bson.ObjectID, workflows []models.Workflow) error
	FindWorkflowByProjectID(ctx context.Context, projectID bson.ObjectID) ([]models.Workflow, error)
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

type CreateProjectMemberRequest struct {
	UserID   bson.ObjectID
	Position string
	Role     models.ProjectMemberRole
}

type SearchProjectMemberRequest struct {
	ProjectID         bson.ObjectID
	Keyword           string
	PaginationRequest PaginationRequest
}
