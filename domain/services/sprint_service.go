package services

import (
	"context"
	"fmt"

	"github.com/cnc-csku/task-nexus-go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SprintService interface {
	Create(ctx context.Context, req *requests.CreateSprintRequest, userID string) (*responses.CreateSprintResponse, *errutils.Error)
	GetByID(ctx context.Context, req *requests.GetSprintByIDRequest) (*models.Sprint, *errutils.Error)
	Edit(ctx context.Context, req *requests.EditSprintRequest, userID string) (*responses.EditSprintResponse, *errutils.Error)
	ListByProjectID(ctx context.Context, req *requests.ListSprintByProjectIDPathParam) ([]models.Sprint, *errutils.Error)
}

type sprintServiceImpl struct {
	sprintRepo  repositories.SprintRepository
	projectRepo repositories.ProjectRepository
}

func NewSprintService(
	sprintRepo repositories.SprintRepository,
	projectRepo repositories.ProjectRepository,
) SprintService {
	return &sprintServiceImpl{
		sprintRepo:  sprintRepo,
		projectRepo: projectRepo,
	}
}

func (s *sprintServiceImpl) Create(ctx context.Context, req *requests.CreateSprintRequest, userID string) (*responses.CreateSprintResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}
	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	project, err := s.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest).WithDebugMessage("project not found")
	}

	// should be in transaction, to be implemented
	sprint := &repositories.CreateSprintRequest{
		ProjectID: bsonProjectID,
		Title:     fmt.Sprintf("%s Sprint %d", project.Name, project.SprintRunningNumber),
		CreatedBy: bsonUserID,
	}

	createdSprint, err := s.sprintRepo.Create(ctx, sprint)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	err = s.projectRepo.IncrementSprintRunningNumber(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return &responses.CreateSprintResponse{
		ID:        createdSprint.ID.Hex(),
		ProjectID: createdSprint.ProjectID.Hex(),
		Title:     createdSprint.Title,
		CreatedAt: createdSprint.CreatedAt,
		CreatedBy: createdSprint.CreatedBy.Hex(),
	}, nil
}

func (s *sprintServiceImpl) GetByID(ctx context.Context, req *requests.GetSprintByIDRequest) (*models.Sprint, *errutils.Error) {
	bsonSprintID, err := bson.ObjectIDFromHex(req.SprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	sprint, err := s.sprintRepo.FindByID(ctx, bsonSprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if sprint == nil {
		return nil, errutils.NewError(exceptions.ErrSprintNotFound, errutils.NotFound).WithDebugMessage("sprint not found")
	}

	return sprint, nil
}

func (s *sprintServiceImpl) Edit(ctx context.Context, req *requests.EditSprintRequest, userID string) (*responses.EditSprintResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}
	bsonSprintID, err := bson.ObjectIDFromHex(req.SprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	sprint, err := s.sprintRepo.FindByID(ctx, bsonSprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if sprint == nil {
		return nil, errutils.NewError(exceptions.ErrSprintNotFound, errutils.NotFound).WithDebugMessage("sprint not found")
	}

	var (
		startDate = req.StartDate
		endDate   = req.EndDate
	)
	if req.Duration != nil {
		if req.StartDate == nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage("start date is required")
		}
		computedEndDate := req.StartDate.AddDate(0, 0, int(*req.Duration))
		endDate = &computedEndDate
	}

	sprintUpdateRequest := &repositories.UpdateSprintRequest{
		ID:         bsonSprintID,
		Title:      req.Title,
		SprintGoal: req.SprintGoal,
		StartDate:  startDate,
		EndDate:    endDate,
		UpdatedBy:  bsonUserID,
	}

	err = s.sprintRepo.Update(ctx, sprintUpdateRequest)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return &responses.EditSprintResponse{
		Message: "Sprint updated successfully",
	}, nil
}

func (s *sprintServiceImpl) ListByProjectID(ctx context.Context, req *requests.ListSprintByProjectIDPathParam) ([]models.Sprint, *errutils.Error) {
	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	project, err := s.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest).WithDebugMessage("project not found")
	}

	sprints, err := s.sprintRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return sprints, nil
}
