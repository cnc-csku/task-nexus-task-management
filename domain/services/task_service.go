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

type TaskService interface {
	Create(ctx context.Context, req *requests.CreateTaskRequest, userID string) (*models.Task, *errutils.Error)
	GetTaskDetail(ctx context.Context, req *requests.GetTaskDetailPathParam, userId string) (*responses.GetTaskDetailResponse, *errutils.Error)
}

type taskServiceImpl struct {
	taskRepo          repositories.TaskRepository
	projectRepo       repositories.ProjectRepository
	projectMemberRepo repositories.ProjectMemberRepository
	sprintRepo        repositories.SprintRepository
	taskCommentRepo   repositories.TaskCommentRepository
	userRepo          repositories.UserRepository
}

func NewTaskService(
	taskRepo repositories.TaskRepository,
	projectRepo repositories.ProjectRepository,
	projectMemberRepo repositories.ProjectMemberRepository,
	sprintRepo repositories.SprintRepository,
	taskCommentRepo repositories.TaskCommentRepository,
	userRepo repositories.UserRepository,
) TaskService {
	return &taskServiceImpl{
		taskRepo:          taskRepo,
		projectRepo:       projectRepo,
		projectMemberRepo: projectMemberRepo,
		sprintRepo:        sprintRepo,
		taskCommentRepo:   taskCommentRepo,
		userRepo:          userRepo,
	}
}

func (s *taskServiceImpl) Create(ctx context.Context, req *requests.CreateTaskRequest, userID string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	var bsonSprintID *bson.ObjectID
	if req.SprintID != nil {
		sprintID, err := bson.ObjectIDFromHex(*req.SprintID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
		}
		bsonSprintID = &sprintID
	}

	// Check if task type is valid
	if !models.TaskType(req.Type).IsValid() {
		return nil, errutils.NewError(exceptions.ErrInvalidTaskType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid task type: %s", req.Type))
	}

	var parentID *string
	if req.ParentID != nil {
		parentTask, err := s.taskRepo.FindByTaskID(ctx, *req.ParentID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if parentTask == nil {
			return nil, errutils.NewError(exceptions.ErrParentTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task not found: %s", *req.ParentID))
		}

		if serviceErr := validateParentTaskType(req.Type, parentTask.Type); serviceErr != nil {
			return nil, serviceErr
		}

		parentID = req.ParentID
	}

	project, err := s.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest)
	}

	var taskSprint *models.TaskSprint
	if bsonSprintID != nil {
		sprint, err := s.sprintRepo.FindByID(ctx, *bsonSprintID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if sprint == nil {
			return nil, errutils.NewError(exceptions.ErrSprintNotFound, errutils.BadRequest)
		}

		taskSprint = &models.TaskSprint{
			CurrentSprintID: *bsonSprintID,
		}
	}

	var defaultWorkflow *models.Workflow
	for _, workflow := range project.Workflows {
		if workflow.IsDefault {
			defaultWorkflow = &workflow
			break
		}
	}
	if defaultWorkflow == nil {
		return nil, errutils.NewError(exceptions.ErrDefaultWorkflowNotFound, errutils.InternalServerError)
	}

	task, err := s.taskRepo.Create(ctx, &repositories.CreateTaskRequest{
		TaskID:      fmt.Sprintf("%s-%d", project.ProjectPrefix, project.TaskRunningNumber),
		ProjectID:   bsonProjectID,
		Title:       req.Title,
		Description: req.Description,
		ParentID:    parentID,
		Type:        models.TaskType(req.Type),
		Status:      defaultWorkflow.Status,
		Sprint:      taskSprint,
		CreatedBy:   bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	err = s.projectRepo.IncrementTaskRunningNumber(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return task, nil
}

func validateParentTaskType(taskType string, parentTaskType models.TaskType) *errutils.Error {
	switch taskType {
	case models.TaskTypeEpic.String():
		return nil
	case models.TaskTypeStory.String(), models.TaskTypeTask.String(), models.TaskTypeBug.String():
		if parentTaskType != models.TaskTypeEpic {
			return errutils.NewError(exceptions.ErrInvalidParentTaskType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task type is not valid: %s", parentTaskType))
		}
	case models.TaskTypeSubTask.String():
		if !array.ContainAny(
			[]string{parentTaskType.String()},
			[]string{
				models.TaskTypeStory.String(),
				models.TaskTypeTask.String(),
				models.TaskTypeBug.String(),
			},
		) {
			return errutils.NewError(exceptions.ErrInvalidParentTaskType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task type is not valid: %s", parentTaskType))
		}
	}
	return nil
}

func (s *taskServiceImpl) GetTaskDetail(ctx context.Context, req *requests.GetTaskDetailPathParam, userId string) (*responses.GetTaskDetailResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskID(ctx, req.TaskID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if task == nil {
		return nil, errutils.NewError(exceptions.ErrTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Task not found: %s", req.TaskID))
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, task.ProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	creator, err := s.userRepo.FindByID(ctx, task.CreatedBy)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if creator == nil {
		return nil, errutils.NewError(exceptions.ErrUserNotFound, errutils.InternalServerError).WithDebugMessage(fmt.Sprintf("User not found: %s", task.CreatedBy.Hex()))
	}

	updater, err := s.userRepo.FindByID(ctx, task.UpdatedBy)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if updater == nil {
		return nil, errutils.NewError(exceptions.ErrUserNotFound, errutils.InternalServerError).WithDebugMessage(fmt.Sprintf("User not found: %s", task.UpdatedBy.Hex()))
	}

	comments, err := s.taskCommentRepo.FindByTaskID(ctx, req.TaskID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	commentsUserIDs := extractUserIDsFromComments(comments)

	commentsUsers, err := s.userRepo.FindByIDs(ctx, commentsUserIDs)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	commentsUsersMap := mapUsersByID(commentsUsers)

	return &responses.GetTaskDetailResponse{
		ID:                 task.ID.Hex(),
		TaskID:             task.TaskID,
		ProjectID:          task.ProjectID.Hex(),
		Title:              task.Title,
		Description:        task.Description,
		ParentID:           task.ParentID,
		Type:               task.Type,
		Status:             task.Status,
		Priority:           task.Priority,
		Approval:           task.Approval,
		Assignee:           task.Assignee,
		Sprint:             task.Sprint,
		CreatedAt:          task.CreatedAt,
		CreatedBy:          task.CreatedBy.Hex(),
		CreatorDisplayName: creator.DisplayName,
		UpdatedAt:          task.UpdatedAt,
		UpdatedBy:          task.UpdatedBy.Hex(),
		UpdaterDisplayName: updater.DisplayName,
		TaskComments:       buildTaskComments(comments, commentsUsersMap),
	}, nil
}

func extractUserIDsFromComments(comments []*models.TaskComment) []bson.ObjectID {
	userIDs := make([]bson.ObjectID, 0, len(comments))
	for _, comment := range comments {
		userIDs = append(userIDs, comment.UserID)
	}
	return userIDs
}

func mapUsersByID(users []models.User) map[string]string {
	userMap := make(map[string]string, len(users))
	for _, user := range users {
		userMap[user.ID.Hex()] = user.DisplayName
	}
	return userMap
}

func buildTaskComments(comments []*models.TaskComment, userMap map[string]string) []responses.GetTaskDetailResponseTaskComment {
	taskComments := make([]responses.GetTaskDetailResponseTaskComment, 0, len(comments))
	for _, comment := range comments {
		taskComments = append(taskComments, responses.GetTaskDetailResponseTaskComment{
			ID:              comment.ID.Hex(),
			Content:         comment.Content,
			UserID:          comment.UserID.Hex(),
			UserDisplayName: userMap[comment.UserID.Hex()],
			TaskID:          comment.TaskID,
			CreatedAt:       comment.CreatedAt,
			UpdatedAt:       comment.UpdatedAt,
		})
	}
	return taskComments
}
