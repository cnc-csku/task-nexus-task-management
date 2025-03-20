package services

import (
	"context"

	"github.com/cnc-csku/task-nexus-go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ReportService interface {
	GetStatusOverview(ctx context.Context, req *requests.GetStatusOverviewRequest, userID string) (*responses.GetStatusOverviewResponse, *errutils.Error)
	// GetPriorityOverview(ctx context.Context, req *requests.GetPriorityOverviewRequest, userID string) (*responses.GetPriorityOverviewResponse, *errutils.Error)
}

type reportServiceImpl struct {
	projectRepo   repositories.ProjectRepository
	projectMember repositories.ProjectMemberRepository
	taskRepo      repositories.TaskRepository
}

func NewReportService(
	projectRepo repositories.ProjectRepository,
	projectMember repositories.ProjectMemberRepository,
	taskRepo repositories.TaskRepository,
) ReportService {
	return &reportServiceImpl{
		projectRepo:   projectRepo,
		projectMember: projectMember,
		taskRepo:      taskRepo,
	}
}

func (s *reportServiceImpl) GetStatusOverview(ctx context.Context, req *requests.GetStatusOverviewRequest, userID string) (*responses.GetStatusOverviewResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	project, err := s.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest).WithDebugMessage("project not found")
	}

	member, err := s.projectMember.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.Forbidden).WithDebugMessage("permission denied")
	}

	tasks, err := s.taskRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if len(tasks) == 0 {
		return &responses.GetStatusOverviewResponse{
			Statuses:   []responses.GetStatusOverviewResponseStatuses{},
			TotalCount: 0,
		}, nil
	}

	statuses := make(map[string]int)
	for _, task := range tasks {
		statuses[task.Status]++
	}

	statusOverview := make([]responses.GetStatusOverviewResponseStatuses, 0)
	for status, count := range statuses {
		statusOverview = append(statusOverview, responses.GetStatusOverviewResponseStatuses{
			Status:  status,
			Count:   count,
			Percent: float64(count) / float64(len(tasks)) * 100,
		})
	}

	return &responses.GetStatusOverviewResponse{
		Statuses:   statusOverview,
		TotalCount: len(tasks),
	}, nil
}
