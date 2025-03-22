package services

import (
	"context"
	"fmt"

	"github.com/cnc-csku/task-nexus-go-lib/utils/array"
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
	GetByID(ctx context.Context, req *requests.GetSprintByIDRequest, userID string) (*models.Sprint, *errutils.Error)
	Edit(ctx context.Context, req *requests.EditSprintRequest, userID string) (*responses.EditSprintResponse, *errutils.Error)
	List(ctx context.Context, req *requests.ListSprintPathParam, userID string) ([]models.Sprint, *errutils.Error)
	CompleteSprint(ctx context.Context, req *requests.CompleteSprintRequest, userID string) (*models.Sprint, *errutils.Error)
	UpdateStatus(ctx context.Context, req *requests.UpdateSprintStatusRequest, userID string) (*models.Sprint, *errutils.Error)
	Delete(ctx context.Context, req *requests.DeleteSprintRequest, userID string) (*responses.DeleteSprintResponse, *errutils.Error)
}

type sprintServiceImpl struct {
	sprintRepo        repositories.SprintRepository
	projectRepo       repositories.ProjectRepository
	projectMemberRepo repositories.ProjectMemberRepository
	taskRepo          repositories.TaskRepository
}

func NewSprintService(
	sprintRepo repositories.SprintRepository,
	projectRepo repositories.ProjectRepository,
	projectMemberRepo repositories.ProjectMemberRepository,
	taskRepo repositories.TaskRepository,
) SprintService {
	return &sprintServiceImpl{
		sprintRepo:        sprintRepo,
		projectRepo:       projectRepo,
		projectMemberRepo: projectMemberRepo,
		taskRepo:          taskRepo,
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
		Status:    models.SprintStatusCreated,
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
		Status:    createdSprint.Status.String(),
		CreatedAt: createdSprint.CreatedAt,
		CreatedBy: createdSprint.CreatedBy.Hex(),
	}, nil
}

func (s *sprintServiceImpl) GetByID(ctx context.Context, req *requests.GetSprintByIDRequest, userID string) (*models.Sprint, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonSprintID, err := bson.ObjectIDFromHex(req.SprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("user is not member of project")
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

func (s *sprintServiceImpl) List(ctx context.Context, req *requests.ListSprintPathParam, userID string) ([]models.Sprint, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	var statuses []models.SprintStatus
	for _, status := range req.Statuses {
		if !models.SprintStatus(status).IsValid() {
			return nil, errutils.NewError(exceptions.ErrInvalidSprintStatus, errutils.BadRequest).WithDebugMessage("invalid sprint status")
		}

		statuses = append(statuses, models.SprintStatus(status))
	}

	project, err := s.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest).WithDebugMessage("project not found")
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("user is not member of project")
	}

	sprints, err := s.sprintRepo.List(ctx, &repositories.ListSprintFilter{
		ProjectID: bsonProjectID,
		IsActive:  req.IsActive,
		Statuses:  statuses,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return sprints, nil
}

func (s *sprintServiceImpl) CompleteSprint(ctx context.Context, req *requests.CompleteSprintRequest, userID string) (*models.Sprint, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonCurrentSprintID, err := bson.ObjectIDFromHex(req.CurrentSprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	project, err := s.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Project not found: %s", req.ProjectID))
	}

	currentSprint, err := s.sprintRepo.FindByID(ctx, bsonCurrentSprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if currentSprint == nil {
		return nil, errutils.NewError(exceptions.ErrSprintNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Sprint not found: %s", req.CurrentSprintID))
	}

	tasks, err := s.taskRepo.FindByCurrentSprintID(ctx, bsonCurrentSprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	var isDoneStatus []string
	for _, workflow := range project.Workflows {
		if workflow.IsDone {
			isDoneStatus = append(isDoneStatus, workflow.Status)
			break
		}
	}

	notDoneTasks := make([]string, 0, len(tasks))
	for _, task := range tasks {
		if !array.ContainAny(isDoneStatus, []string{task.Status}) {
			notDoneTasks = append(notDoneTasks, task.TaskRef)
		}
	}

	if len(notDoneTasks) > 0 {
		return nil, errutils.NewError(
			exceptions.ErrNotAllTasksIsDone, errutils.BadRequest,
		).WithDebugMessage(
			fmt.Sprintf("Not all tasks are done: %v", notDoneTasks),
		).WithFields(
			notDoneTasks...,
		)
	}

	updatedSprint, err := s.sprintRepo.UpdateStatus(ctx, &repositories.UpdateSprintStatusRequest{
		ID:        bsonCurrentSprintID,
		Status:    models.SprintStatusCompleted,
		UpdatedBy: bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return updatedSprint, nil
}

func (s *sprintServiceImpl) UpdateStatus(ctx context.Context, req *requests.UpdateSprintStatusRequest, userID string) (*models.Sprint, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonSprintID, err := bson.ObjectIDFromHex(req.SprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	if !models.SprintStatus(req.Status).IsValid() {
		return nil, errutils.NewError(exceptions.ErrInvalidSprintStatus, errutils.BadRequest).WithDebugMessage("invalid sprint status")
	}

	project, err := s.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest).WithDebugMessage("project not found")
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("user is not member of project")
	}

	sprint, err := s.sprintRepo.FindByID(ctx, bsonSprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if sprint == nil {
		return nil, errutils.NewError(exceptions.ErrSprintNotFound, errutils.NotFound).WithDebugMessage("sprint not found")
	}

	sprint, err = s.sprintRepo.UpdateStatus(ctx, &repositories.UpdateSprintStatusRequest{
		ID:     bsonSprintID,
		Status: models.SprintStatus(req.Status),
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return sprint, nil
}

func (s *sprintServiceImpl) Delete(ctx context.Context, req *requests.DeleteSprintRequest, userID string) (*responses.DeleteSprintResponse, *errutils.Error) {
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

	// should be in transaction, to be implemented
	tasksWithCurrentSprintID, err := s.taskRepo.FindByCurrentSprintID(ctx, bsonSprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if tasksWithCurrentSprintID != nil {
		for _, task := range tasksWithCurrentSprintID {
			_, err = s.taskRepo.UpdateCurrentSprintID(ctx, &repositories.UpdateTaskCurrentSprintIDRequest{
				ID:              task.ID,
				CurrentSprintID: nil,
				UpdatedBy:       bsonUserID,
			})
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}
		}
	}

	tasksWithPreviousSprintID, err := s.taskRepo.FindByPreviousSprintID(ctx, bsonSprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if tasksWithPreviousSprintID != nil {
		for _, task := range tasksWithPreviousSprintID {
			prevSprintIDs, err := array.RemoveOne(task.Sprint.PreviousSprintIDs, bsonSprintID)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}

			_, err = s.taskRepo.UpdatePreviousSprintIDs(ctx, &repositories.UpdateTaskPreviousSprintIDsRequest{
				ID:                task.ID,
				PreviousSprintIDs: prevSprintIDs,
				UpdatedBy:         bsonUserID,
			})
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}
		}
	}

	err = s.sprintRepo.Delete(ctx, bsonSprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return &responses.DeleteSprintResponse{
		Message: "Sprint deleted successfully",
	}, nil
}
