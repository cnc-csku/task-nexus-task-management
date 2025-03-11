package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cnc-csku/task-nexus-go-lib/utils/array"
	"github.com/cnc-csku/task-nexus-go-lib/utils/conv"
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
	ListEpicTasks(ctx context.Context, req *requests.ListEpicTasksPathParam, userId string) ([]*models.Task, *errutils.Error)
	SearchTask(ctx context.Context, req *requests.SearchTaskParams, userId string) ([]responses.SearchTaskResponse, *errutils.Error)
	UpdateDetail(ctx context.Context, req *requests.UpdateTaskDetailRequest, userId string) (*models.Task, *errutils.Error)
	UpdateStatus(ctx context.Context, req *requests.UpdateTaskStatusRequest, userId string) (*models.Task, *errutils.Error)
	UpdateApprovals(ctx context.Context, req *requests.UpdateTaskApprovalsRequest, userId string) (*models.Task, *errutils.Error)
	ApproveTask(ctx context.Context, req *requests.ApproveTaskRequest, userId string) (*models.Task, *errutils.Error)
	UpdateAssignees(ctx context.Context, req *requests.UpdateTaskAssigneesRequest, userId string) (*models.Task, *errutils.Error)
	UpdateSprint(ctx context.Context, req *requests.UpdateTaskSprintRequest, userId string) (*models.Task, *errutils.Error)
	UpdateAttributes(ctx context.Context, req *requests.UpdateTaskAttributesRequest, userId string) (*models.Task, *errutils.Error)
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

	var nullableBsonTaskParentID *bson.ObjectID
	if req.ParentID != nil {
		bsonTaskParentID, err := bson.ObjectIDFromHex(*req.ParentID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
		}

		parentTask, err := s.taskRepo.FindByID(ctx, bsonTaskParentID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if parentTask == nil {
			return nil, errutils.NewError(exceptions.ErrParentTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task not found: %s", *req.ParentID))
		}

		if serviceErr := validateParentTaskType(req.Type, parentTask.Type); serviceErr != nil {
			return nil, serviceErr
		}

		_, err = s.taskRepo.UpdateHasChildren(ctx, &repositories.UpdateTaskHasChildrenRequest{
			ID:          bsonTaskParentID,
			HasChildren: true,
		})
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		nullableBsonTaskParentID = &bsonTaskParentID
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

	var defaultWorkflow *models.ProjectWorkflow
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
		TaskRef:     fmt.Sprintf("%s-%d", project.ProjectPrefix, project.TaskRunningNumber),
		ProjectID:   bsonProjectID,
		Title:       req.Title,
		Description: req.Description,
		ParentID:    nullableBsonTaskParentID,
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

func validateParentTaskType(childTaskType string, parentTaskType models.TaskType) *errutils.Error {
	switch childTaskType {
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

	task, err := s.taskRepo.FindByTaskRef(ctx, req.TaskRef)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if task == nil {
		return nil, errutils.NewError(exceptions.ErrTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Task not found: %s", req.TaskRef))
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

	comments, err := s.taskCommentRepo.FindByTaskID(ctx, task.ID)
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
		TaskRef:            task.TaskRef,
		ProjectID:          task.ProjectID.Hex(),
		Title:              task.Title,
		Description:        task.Description,
		ParentID:           task.ParentID,
		Type:               task.Type,
		Status:             task.Status,
		Priority:           task.Priority,
		Approvals:          task.Approvals,
		Assignees:          task.Assignees,
		ChildrenPoint:      task.ChildrenPoint,
		HasChildren:        task.HasChildren,
		Sprint:             task.Sprint,
		Attributes:         task.Attributes,
		StartDate:          task.StartDate,
		DueDate:            task.DueDate,
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
			TaskID:          comment.TaskID.Hex(),
			CreatedAt:       comment.CreatedAt,
			UpdatedAt:       comment.UpdatedAt,
		})
	}
	return taskComments
}

func (s *taskServiceImpl) ListEpicTasks(ctx context.Context, req *requests.ListEpicTasksPathParam, userId string) ([]*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
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

	tasks, err := s.taskRepo.FindByProjectIDAndType(ctx, bsonProjectID, models.TaskTypeEpic)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return tasks, nil
}

func (s *taskServiceImpl) SearchTask(ctx context.Context, req *requests.SearchTaskParams, userId string) ([]responses.SearchTaskResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	var bsonSprintID *bson.ObjectID
	if req.SprintID != nil {
		sprintID, err := bson.ObjectIDFromHex(*req.SprintID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
		}
		bsonSprintID = &sprintID
	}

	var bsonParentID *bson.ObjectID
	if req.EpicTaskID != nil {
		parentID, err := bson.ObjectIDFromHex(*req.EpicTaskID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
		}
		bsonParentID = &parentID

		parentTask, err := s.taskRepo.FindByID(ctx, parentID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if parentTask == nil {
			return nil, errutils.NewError(exceptions.ErrParentTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task not found: %s", *req.EpicTaskID))
		} else if parentTask.Type != models.TaskTypeEpic {
			return nil, errutils.NewError(exceptions.ErrInvalidParentTaskType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task type is not valid: %s", parentTask.Type))
		}
	}

	userIDs := make([]bson.ObjectID, 0, len(req.UserIDs))
	for _, userID := range req.UserIDs {
		bsonUserID, err := bson.ObjectIDFromHex(userID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
		}
		userIDs = append(userIDs, bsonUserID)
	}

	project, err := s.projectRepo.FindByProjectID(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest)
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	tasks, err := s.taskRepo.Search(ctx, &repositories.SearchTaskRequest{
		ProjectID:      bsonProjectID,
		TaskTypes:      []models.TaskType{models.TaskTypeStory, models.TaskTypeTask, models.TaskTypeBug},
		SprintID:       bsonSprintID,
		EpicTaskID:     bsonParentID,
		UserIDs:        userIDs,
		Positions:      req.Positions,
		Statuses:       req.Statuses,
		IsDoneStatuses: getDoneStatuses(project),
		SearchKeyword:  req.SearchKeyword,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if len(tasks) == 0 {
		return []responses.SearchTaskResponse{}, nil
	}

	parentTasksMap, serviceErr := getParentTasksMap(ctx, s.taskRepo, tasks)
	if serviceErr != nil {
		return nil, serviceErr
	}

	response := make([]responses.SearchTaskResponse, 0, len(tasks))
	for _, task := range tasks {
		response = append(response, responses.SearchTaskResponse{
			ID:       task.ID.Hex(),
			TaskRef:  task.TaskRef,
			Title:    task.Title,
			ParentID: conv.BsonObjectIDPtrToStringPtr(task.ParentID),
			ParentTitle: func() *string {
				if task.ParentID == nil {
					return nil
				}
				return parentTasksMap[task.ParentID.Hex()]
			}(),
			Type:          task.Type.String(),
			Status:        task.Status,
			Assignees:     task.Assignees,
			ChildrenPoint: task.ChildrenPoint,
			HasChildren:   task.HasChildren,
			Sprint:        task.Sprint,
		})
	}

	return response, nil
}

func getDoneStatuses(project *models.Project) []string {
	isDoneStatuses := make([]string, 0)
	for _, workflow := range project.Workflows {
		if workflow.IsDone {
			isDoneStatuses = append(isDoneStatuses, workflow.Status)
		}
	}
	return isDoneStatuses
}

func getParentTasksMap(ctx context.Context, taskRepo repositories.TaskRepository, tasks []*models.Task) (map[string]*string, *errutils.Error) {
	parentTaskIDs := make([]bson.ObjectID, 0, len(tasks))
	for _, task := range tasks {
		if task.ParentID != nil {
			parentTaskIDs = append(parentTaskIDs, *task.ParentID)
		}
	}

	parentTasks, err := taskRepo.FindByIDs(ctx, parentTaskIDs)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	parentTasksMap := make(map[string]*string, len(parentTasks))
	for _, parentTask := range parentTasks {
		parentTasksMap[parentTask.ID.Hex()] = &parentTask.Title
	}

	return parentTasksMap, nil
}

func (s *taskServiceImpl) UpdateDetail(ctx context.Context, req *requests.UpdateTaskDetailRequest, userId string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	// Check if task priority is valid
	if !models.TaskPriority(req.Priority).IsValid() {
		return nil, errutils.NewError(exceptions.ErrInvalidTaskPriority, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid task priority: %s", req.Priority))
	}

	task, err := s.taskRepo.FindByTaskRef(ctx, req.TaskRef)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if task == nil {
		return nil, errutils.NewError(exceptions.ErrTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Task not found: %s", req.TaskRef))
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, task.ProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	var (
		nullableBsonTaskParentID = task.ParentID
		isParentTaskChanged      = task.ParentID == nil && req.ParentID != nil || req.ParentID != nil && task.ParentID != nil && *req.ParentID != task.ParentID.Hex()
		isParentTaskRemoved      = req.ParentID == nil && task.ParentID != nil
	)
	if isParentTaskChanged {
		bsonNewTaskParentID, err := bson.ObjectIDFromHex(*req.ParentID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
		}

		newParentTask, err := s.taskRepo.FindByID(ctx, bsonNewTaskParentID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if newParentTask == nil {
			return nil, errutils.NewError(exceptions.ErrParentTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task not found: %s", *req.ParentID))
		}

		serviceErr := validateParentTaskType(task.Type.String(), newParentTask.Type)
		if serviceErr != nil {
			return nil, serviceErr
		}

		serviceErr = updatePreviousParentTask(ctx, s.taskRepo, task)
		if serviceErr != nil {
			return nil, serviceErr
		}

		if !newParentTask.HasChildren {
			_, err = s.taskRepo.UpdateHasChildren(ctx, &repositories.UpdateTaskHasChildrenRequest{
				ID:          bsonNewTaskParentID,
				HasChildren: true,
			})
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}
		}

		if task.Type == models.TaskTypeSubTask {
			var totalPoint int
			for _, assignee := range task.Assignees {
				if assignee.Point != nil {
					totalPoint += *assignee.Point
				}
			}

			_, err = s.taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
				ID:            bsonNewTaskParentID,
				ChildrenPoint: newParentTask.ChildrenPoint + totalPoint,
			})
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}
		} else if array.ContainAny(
			[]string{task.Type.String()},
			[]string{
				models.TaskTypeStory.String(),
				models.TaskTypeTask.String(),
				models.TaskTypeBug.String(),
			}) {
			_, err = s.taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
				ID:            bsonNewTaskParentID,
				ChildrenPoint: newParentTask.ChildrenPoint + task.ChildrenPoint,
			})
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}
		}

		nullableBsonTaskParentID = &bsonNewTaskParentID
	} else if isParentTaskRemoved {
		serviceErr := updatePreviousParentTask(ctx, s.taskRepo, task)
		if serviceErr != nil {
			return nil, serviceErr
		}

		nullableBsonTaskParentID = nil
	}

	updatedTask, err := s.taskRepo.UpdateDetail(ctx, &repositories.UpdateTaskDetailRequest{
		ID:          task.ID,
		Title:       req.Title,
		Description: req.Description,
		ParentID:    nullableBsonTaskParentID,
		Priority:    req.Priority,
		StartDate:   req.StartDate,
		DueDate:     req.DueDate,
		UpdatedBy:   bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return updatedTask, nil
}

func updatePreviousParentTask(
	ctx context.Context,
	taskRepo repositories.TaskRepository,
	updatedTask *models.Task,
) *errutils.Error {
	if updatedTask.ParentID == nil {
		return nil
	}

	previousParentTask, err := taskRepo.FindByID(ctx, *updatedTask.ParentID)
	if err != nil {
		return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if previousParentTask == nil {
		return errutils.NewError(exceptions.ErrParentTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task not found: %s", *updatedTask.ParentID))
	}

	childrenOfPreviousParentTasks, err := taskRepo.FindByParentID(ctx, previousParentTask.ID)
	if err != nil {
		return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	if len(childrenOfPreviousParentTasks) == 0 {
		_, err := taskRepo.UpdateHasChildren(ctx, &repositories.UpdateTaskHasChildrenRequest{
			ID:          previousParentTask.ID,
			HasChildren: false,
		})
		if err != nil {
			return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}
	}

	serviceErr := updateChildrenPointOfPreviousParentTask(ctx, taskRepo, updatedTask, previousParentTask)
	if serviceErr != nil {
		return serviceErr
	}

	return nil
}

func updateChildrenPointOfPreviousParentTask(
	ctx context.Context,
	taskRepo repositories.TaskRepository,
	updatedTask, previousParentTask *models.Task,
) *errutils.Error {
	if previousParentTask.Type == models.TaskTypeEpic {
		_, err := taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
			ID:            previousParentTask.ID,
			ChildrenPoint: previousParentTask.ChildrenPoint - updatedTask.ChildrenPoint,
		})
		if err != nil {
			return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}
	} else if array.ContainAny(
		[]string{previousParentTask.Type.String()},
		[]string{
			models.TaskTypeStory.String(),
			models.TaskTypeTask.String(),
			models.TaskTypeBug.String(),
		},
	) {
		// Current previous parent task is a story, task, or bug
		var totalPoint int
		for _, assignee := range updatedTask.Assignees {
			if assignee.Point != nil {
				totalPoint += *assignee.Point
			}
		}

		_, err := taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
			ID:            previousParentTask.ID,
			ChildrenPoint: previousParentTask.ChildrenPoint - totalPoint,
		})
		if err != nil {
			return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		// Update children point of the previous parent task's parent task (EPIC)
		if previousParentTask.ParentID != nil {
			epicParentTask, err := taskRepo.FindByID(ctx, *previousParentTask.ParentID)
			if err != nil {
				return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}

			_, err = taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
				ID:            *previousParentTask.ParentID,
				ChildrenPoint: epicParentTask.ChildrenPoint - totalPoint,
			})
			if err != nil {
				return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}
		}
	}

	return nil
}

func (s *taskServiceImpl) UpdateStatus(ctx context.Context, req *requests.UpdateTaskStatusRequest, userId string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRef(ctx, req.TaskID)
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

	project, err := s.projectRepo.FindByProjectID(ctx, task.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest)
	}

	statuses := make([]string, 0, len(project.Workflows))
	for _, workflow := range project.Workflows {
		statuses = append(statuses, workflow.Status)
	}

	if !array.ContainAny(statuses, []string{req.Status}) {
		return nil, errutils.NewError(exceptions.ErrInvalidTaskStatus, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid task status: %s", req.Status))
	}

	updatedTask, err := s.taskRepo.UpdateStatus(ctx, &repositories.UpdateTaskStatusRequest{
		ID:        task.ID,
		Status:    req.Status,
		UpdatedBy: bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return updatedTask, nil
}

func (s *taskServiceImpl) UpdateApprovals(ctx context.Context, req *requests.UpdateTaskApprovalsRequest, userID string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRef(ctx, req.TaskRef)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if task == nil {
		return nil, errutils.NewError(exceptions.ErrTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Task not found: %s", req.TaskRef))
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, task.ProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	approvals := make([]repositories.UpdateTaskApprovalsRequestApproval, 0, len(req.ApprovalUserIDs))
	for _, userID := range req.ApprovalUserIDs {
		bsonApprovalUserID, err := bson.ObjectIDFromHex(userID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
		}

		approvals = append(approvals, repositories.UpdateTaskApprovalsRequestApproval{
			UserID: bsonApprovalUserID,
		})
	}

	updatedTask, err := s.taskRepo.UpdateApprovals(ctx, &repositories.UpdateTaskApprovalsRequest{
		ID:        task.ID,
		Approval:  approvals,
		UpdatedBy: bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return updatedTask, nil
}

func (s *taskServiceImpl) ApproveTask(ctx context.Context, req *requests.ApproveTaskRequest, userId string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRef(ctx, req.TaskRef)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if task == nil {
		return nil, errutils.NewError(exceptions.ErrTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Task not found: %s", req.TaskRef))
	}

	approvalUserIDs := make([]string, 0, len(task.Approvals))
	for _, approval := range task.Approvals {
		approvalUserIDs = append(approvalUserIDs, approval.UserID.Hex())
	}

	if !array.ContainAny(approvalUserIDs, []string{bsonUserID.Hex()}) {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not in the approval list")
	}

	updatedTask, err := s.taskRepo.ApproveTask(ctx, &repositories.ApproveTaskRequest{
		ID:     task.ID,
		Reason: req.Reason,
		UserID: bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return updatedTask, nil
}

func (s *taskServiceImpl) UpdateAssignees(ctx context.Context, req *requests.UpdateTaskAssigneesRequest, userId string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRef(ctx, req.TaskRef)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if task == nil {
		return nil, errutils.NewError(exceptions.ErrTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Task not found: %s", req.TaskRef))
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, task.ProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	var (
		currentTotalPoint int
		assignees         = make([]repositories.UpdateTaskAssigneesRequestAssignee, 0, len(req.Assignees))
	)
	for _, assignee := range task.Assignees {
		if assignee.Point != nil {
			currentTotalPoint += *assignee.Point
		}
	}

	// Modify Assignees, their positions and `point` for SubTask
	if task.Type == models.TaskTypeSubTask {
		for _, assignee := range req.Assignees {
			bsonAssigneeUserID, err := bson.ObjectIDFromHex(assignee.UserId)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
			}

			assignees = append(assignees, repositories.UpdateTaskAssigneesRequestAssignee{
				Position: assignee.Position,
				UserID:   bsonAssigneeUserID,
				Point:    assignee.Point,
			})
		}

		serviceErr := updateParentTaskPoints(ctx, s.taskRepo, task, assignees, currentTotalPoint)
		if serviceErr != nil {
			return nil, serviceErr
		}
	} else {
		// Modify Only Assignees and their positions, not point for Epic, Story, Task, Bug
		for _, assignee := range req.Assignees {
			bsonAssigneeUserID, err := bson.ObjectIDFromHex(assignee.UserId)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
			}

			assignees = append(assignees, repositories.UpdateTaskAssigneesRequestAssignee{
				Position: assignee.Position,
				UserID:   bsonAssigneeUserID,
			})
		}
	}

	updatedTask, err := s.taskRepo.UpdateAssignees(ctx, &repositories.UpdateTaskAssigneesRequest{
		ID:        task.ID,
		Assignees: assignees,
		UpdatedBy: bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return updatedTask, nil
}

func updateParentTaskPoints(
	ctx context.Context,
	taskRepo repositories.TaskRepository,
	task *models.Task,
	assignees []repositories.UpdateTaskAssigneesRequestAssignee,
	currentTotalPoint int,
) *errutils.Error {
	if task.ParentID == nil {
		return nil
	}

	parentTask, err := taskRepo.FindByID(ctx, *task.ParentID)
	if err != nil {
		return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if parentTask == nil {
		return errutils.NewError(exceptions.ErrParentTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task not found: %s", task.ParentID.Hex()))
	}

	var newTotalPoint int
	for _, assignee := range assignees {
		if assignee.Point != nil {
			newTotalPoint += *assignee.Point
		}
	}

	if parentTask.Type == models.TaskTypeEpic {
		_, err := taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
			ID:            *task.ParentID,
			ChildrenPoint: parentTask.ChildrenPoint - currentTotalPoint + newTotalPoint,
		})
		if err != nil {
			return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}
	} else if array.ContainAny(
		[]string{parentTask.Type.String()},
		[]string{
			models.TaskTypeStory.String(),
			models.TaskTypeTask.String(),
			models.TaskTypeBug.String(),
		},
	) {
		// Modify point for Story, Task, Bug
		_, err := taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
			ID:            *task.ParentID,
			ChildrenPoint: parentTask.ChildrenPoint - currentTotalPoint + newTotalPoint,
		})
		if err != nil {
			return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		// Update children point of the parent task's parent task (EPIC)
		if parentTask.ParentID != nil {
			epicParentTask, err := taskRepo.FindByID(ctx, *parentTask.ParentID)
			if err != nil {
				return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}

			_, err = taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
				ID:            *parentTask.ParentID,
				ChildrenPoint: epicParentTask.ChildrenPoint - currentTotalPoint + newTotalPoint,
			})
			if err != nil {
				return errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}
		}
	}

	return nil
}

func (s *taskServiceImpl) UpdateSprint(ctx context.Context, req *requests.UpdateTaskSprintRequest, userId string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRef(ctx, req.TaskRef)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if task == nil {
		return nil, errutils.NewError(exceptions.ErrTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Task not found: %s", req.TaskRef))
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, task.ProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	currentSprintID, err := bson.ObjectIDFromHex(req.CurrentSprintID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	previousSprintIDs := make([]bson.ObjectID, 0, len(req.PreviousSprintIDs))
	for _, previousSprintID := range req.PreviousSprintIDs {
		bsonPreviousSprintID, err := bson.ObjectIDFromHex(previousSprintID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
		}

		previousSprintIDs = append(previousSprintIDs, bsonPreviousSprintID)
	}

	if task.Sprint == nil {

	}

	updatedTask, err := s.taskRepo.UpdateSprint(ctx, &repositories.UpdateTaskSprintRequest{
		ID:                task.ID,
		CurrentSprintID:   currentSprintID,
		PreviousSprintIDs: previousSprintIDs,
		UpdatedBy:         bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return updatedTask, nil
}

func (s *taskServiceImpl) UpdateAttributes(ctx context.Context, req *requests.UpdateTaskAttributesRequest, userId string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRef(ctx, req.TaskRef)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if task == nil {
		return nil, errutils.NewError(exceptions.ErrTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Task not found: %s", req.TaskRef))
	}

	project, err := s.projectRepo.FindByProjectID(ctx, task.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if project == nil {
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Project not found: %s", task.ProjectID.Hex()))
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, task.ProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	attributeTemplateMap := make(map[string]models.ProjectAttributeTemplate)
	for _, attributeTemplate := range project.AttributeTemplates {
		attributeTemplateMap[attributeTemplate.Name] = attributeTemplate
	}

	attributes := make([]models.TaskAttribute, 0, len(req.Attributes))
	for _, attribute := range req.Attributes {
		attributeTemplate, ok := attributeTemplateMap[attribute.Key]
		if !ok {
			return nil, errutils.NewError(exceptions.ErrInvalidAttributeKey, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid attribute key: %s", attribute.Key))
		}

		var value any
		switch attributeTemplate.Type {
		case models.KeyValuePairTypeString:
			value = attribute.Value
		case models.KeyValuePairTypeNumber:
			value, err = strconv.Atoi(attribute.Value)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInvalidAttributeType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid attribute type: %s", attributeTemplate.Type))
			}
		case models.KeyValuePairTypeDate:
			value, err = time.Parse(time.RFC3339, attribute.Value)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInvalidAttributeType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid attribute type: %s", attributeTemplate.Type))
			}
		case models.KeyValuePairTypeBoolean:
			value, err = strconv.ParseBool(attribute.Value)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInvalidAttributeType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid attribute type: %s", attributeTemplate.Type))
			}
		default:
			return nil, errutils.NewError(exceptions.ErrInvalidAttributeType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid attribute type: %s", attributeTemplate.Type))
		}

		attributes = append(attributes, models.TaskAttribute{
			Key:   attribute.Key,
			Value: value,
		})
	}

	updatedTask, err := s.taskRepo.UpdateAttributes(ctx, &repositories.UpdateTaskAttributesRequest{
		ID:         task.ID,
		Attributes: attributes,
		UpdatedBy:  bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return updatedTask, nil
}
