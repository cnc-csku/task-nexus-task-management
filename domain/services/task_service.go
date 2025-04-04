package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cnc-csku/task-nexus-go-lib/utils/array"
	"github.com/cnc-csku/task-nexus-go-lib/utils/conv"
	"github.com/cnc-csku/task-nexus-go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/constant"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
	"github.com/google/generative-ai-go/genai"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskService interface {
	Create(ctx context.Context, req *requests.CreateTaskRequest, userID string) (*models.Task, *errutils.Error)
	GetTaskDetail(ctx context.Context, req *requests.GetTaskDetailPathParam, userId string) (*responses.GetTaskDetailResponse, *errutils.Error)
	GetManyTaskDetail(ctx context.Context, req *requests.GetManyTaskDetailParams, userId string) ([]responses.GetTaskDetailResponse, *errutils.Error)
	ListEpicTasks(ctx context.Context, req *requests.ListEpicTasksPathParam, userId string) ([]*models.Task, *errutils.Error)
	SearchTask(ctx context.Context, req *requests.SearchTaskParams, userId string) ([]responses.SearchTaskResponse, *errutils.Error)
	GetChildrenTasks(ctx context.Context, req *requests.GetChildrenTasksParams, userId string) ([]responses.GetChildrenTasksResponse, *errutils.Error)
	UpdateDetail(ctx context.Context, req *requests.UpdateTaskDetailRequest, userId string) (*models.Task, *errutils.Error)
	UpdateTitle(ctx context.Context, req *requests.UpdateTaskTitleRequest, userId string) (*models.Task, *errutils.Error)
	UpdateParentID(ctx context.Context, req *requests.UpdateTaskParentIdRequest, userId string) (*models.Task, *errutils.Error)
	UpdateType(ctx context.Context, req *requests.UpdateTaskTypeRequest, userId string) (*models.Task, *errutils.Error)
	UpdateStatus(ctx context.Context, req *requests.UpdateTaskStatusRequest, userId string) (*models.Task, *errutils.Error)
	UpdateApprovals(ctx context.Context, req *requests.UpdateTaskApprovalsRequest, userId string) (*models.Task, *errutils.Error)
	ApproveTask(ctx context.Context, req *requests.ApproveTaskRequest, userId string) (*models.Task, *errutils.Error)
	UpdateAssignees(ctx context.Context, req *requests.UpdateTaskAssigneesRequest, userId string) (*models.Task, *errutils.Error)
	UpdateSprint(ctx context.Context, req *requests.UpdateTaskSprintRequest, userId string) (*models.Task, *errutils.Error)
	UpdateAttributes(ctx context.Context, req *requests.UpdateTaskAttributesRequest, userId string) (*models.Task, *errutils.Error)
	GenerateDescription(ctx context.Context, req *requests.GenerateDescriptionRequest, userId string) (*responses.GenerateDescriptionResponse, *errutils.Error)
}

type taskServiceImpl struct {
	taskRepo          repositories.TaskRepository
	projectRepo       repositories.ProjectRepository
	projectMemberRepo repositories.ProjectMemberRepository
	sprintRepo        repositories.SprintRepository
	taskCommentRepo   repositories.TaskCommentRepository
	userRepo          repositories.UserRepository
	geminiRepo        repositories.GeminiRepository
}

func NewTaskService(
	taskRepo repositories.TaskRepository,
	projectRepo repositories.ProjectRepository,
	projectMemberRepo repositories.ProjectMemberRepository,
	sprintRepo repositories.SprintRepository,
	taskCommentRepo repositories.TaskCommentRepository,
	userRepo repositories.UserRepository,
	geminiRepo repositories.GeminiRepository,
) TaskService {
	return &taskServiceImpl{
		taskRepo:          taskRepo,
		projectRepo:       projectRepo,
		projectMemberRepo: projectMemberRepo,
		sprintRepo:        sprintRepo,
		taskCommentRepo:   taskCommentRepo,
		userRepo:          userRepo,
		geminiRepo:        geminiRepo,
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

	if req.DueDate != nil && req.StartDate != nil {
		if req.DueDate.Before(*req.StartDate) {
			return nil, errutils.NewError(exceptions.ErrDueDateBeforeStartDate, errutils.BadRequest).WithDebugMessage("Due date is before start date")
		}
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

	var priority models.TaskPriority
	if req.Priority != nil {
		priority = models.TaskPriority(*req.Priority)
		if !priority.IsValid() {
			return nil, errutils.NewError(exceptions.ErrInvalidTaskPriority, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid task priority: %s", *req.Priority))
		}
	} else {
		priority = models.TaskPriorityMedium
	}

	var (
		bsonSprintID *bson.ObjectID
		startDate    = req.StartDate
		dueDate      = req.DueDate
	)
	if req.SprintID != nil {
		sprintID, err := bson.ObjectIDFromHex(*req.SprintID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
		}

		sprint, err := s.sprintRepo.FindByID(ctx, sprintID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if sprint == nil {
			return nil, errutils.NewError(exceptions.ErrSprintNotFound, errutils.BadRequest)
		}

		bsonSprintID = &sprintID

		if sprint.StartDate != nil && sprint.EndDate != nil {
			startDate = sprint.StartDate
			dueDate = sprint.EndDate
		}
	}

	// Check if task type is valid
	if !models.TaskType(req.Type).IsValid() {
		return nil, errutils.NewError(exceptions.ErrInvalidTaskType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid task type: %s", req.Type))
	}

	var taskSprint = &models.TaskSprint{}
	if bsonSprintID != nil {
		sprint, err := s.sprintRepo.FindByID(ctx, *bsonSprintID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if sprint == nil {
			return nil, errutils.NewError(exceptions.ErrSprintNotFound, errutils.BadRequest)
		}

		taskSprint = &models.TaskSprint{
			CurrentSprintID: bsonSprintID,
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

	assignees := make([]models.TaskAssignee, 0, len(req.Assignees))
	for _, assignee := range req.Assignees {
		var bsonAssigneeUserID *bson.ObjectID
		if assignee.UserID != nil {
			assigneeUserID, err := bson.ObjectIDFromHex(*assignee.UserID)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
			}
			bsonAssigneeUserID = &assigneeUserID
		}

		assignees = append(assignees, models.TaskAssignee{
			UserID:   bsonAssigneeUserID,
			Position: assignee.Position,
			Point:    assignee.Point,
		})
	}

	approvals := make([]models.TaskApproval, 0, len(req.ApprovalUserIDs))
	for _, approvalUserID := range req.ApprovalUserIDs {
		approvalUserID, err := bson.ObjectIDFromHex(approvalUserID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
		}

		approvals = append(approvals, models.TaskApproval{
			UserID: approvalUserID,
		})
	}

	attributeTemplates := make(map[string]models.ProjectAttributeTemplate)
	for _, attributeTemplate := range project.AttributeTemplates {
		attributeTemplates[attributeTemplate.Name] = attributeTemplate
	}

	attributes := make([]models.TaskAttribute, 0, len(req.AdditionalFields))
	for key, value := range req.AdditionalFields {
		if attribute, ok := attributeTemplates[key]; ok {
			if value == nil {
				attributes = append(attributes, models.TaskAttribute{
					Key:   key,
					Value: nil,
				})
				continue
			}

			var val any
			switch attribute.Type {
			case models.KeyValuePairTypeString:
				val, ok = value.(string)
				if !ok {
					return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid attribute value type: %s", key))
				}
			case models.KeyValuePairTypeNumber:
				switch v := value.(type) {
				case float64:
					val = v
				case string:
					val, err = strconv.ParseFloat(v, 64)
					if err != nil {
						return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
					}
				default:
					return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid attribute value type: %s", key))
				}
			case models.KeyValuePairTypeDate:
				val, err = time.Parse(time.RFC3339, value.(string))
				if err != nil {
					return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
				}
			case models.KeyValuePairTypeBoolean:
				switch v := value.(type) {
				case bool:
					val = v
				case string:
					val, err = strconv.ParseBool(v)
					if err != nil {
						return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
					}
				default:
					return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid attribute value type: %s", key))
				}
			default:
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid attribute type: %s", attribute.Type))
			}

			attributes = append(attributes, models.TaskAttribute{
				Key:   key,
				Value: val,
			})
		}
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

		if models.TaskType(req.Type) == models.TaskTypeSubTask {
			var totalPoint int
			for _, assignee := range req.Assignees {
				if assignee.Point != nil {
					totalPoint += *assignee.Point
				}
			}

			_, err := s.taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
				ID:            bsonTaskParentID,
				ChildrenPoint: parentTask.ChildrenPoint + totalPoint,
			})
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}

			if parentTask.ParentID != nil {
				epicParentTask, err := s.taskRepo.FindByID(ctx, *parentTask.ParentID)
				if err != nil {
					return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
				} else if epicParentTask == nil {
					return nil, errutils.NewError(exceptions.ErrParentTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task not found: %s", parentTask.ParentID.Hex()))
				}

				_, err = s.taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
					ID:            *parentTask.ParentID,
					ChildrenPoint: epicParentTask.ChildrenPoint + totalPoint,
				})
				if err != nil {
					return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
				}

				if parentTask.Sprint != nil {
					taskSprint.CurrentSprintID = parentTask.Sprint.CurrentSprintID
				}
			}
		}

		nullableBsonTaskParentID = &bsonTaskParentID
	}

	task, err := s.taskRepo.Create(ctx, &repositories.CreateTaskRequest{
		TaskRef:     fmt.Sprintf("%s-%d", project.ProjectPrefix, project.TaskRunningNumber),
		ProjectID:   bsonProjectID,
		Title:       req.Title,
		Description: req.Description,
		ParentID:    nullableBsonTaskParentID,
		Type:        models.TaskType(req.Type),
		Status:      defaultWorkflow.Status,
		Priority:    priority,
		Sprint:      taskSprint,
		StartDate:   startDate,
		DueDate:     dueDate,
		Assignees:   assignees,
		Approvals:   approvals,
		Attributes:  attributes,
		CreatedBy:   bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	err = s.projectRepo.IncrementTaskRunningNumber(ctx, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	if models.TaskType(req.Type) == models.TaskTypeSubTask {
		if task.ParentID != nil {
			parentTask, err := s.taskRepo.FindByID(ctx, *task.ParentID)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			} else if parentTask == nil {
				return nil, errutils.NewError(exceptions.ErrParentTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task not found: %s", *task.ParentID))
			}

			childrenTasks, err := s.taskRepo.FindByParentID(ctx, parentTask.ID)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}

			serviceErr := updateParentTaskStatusToLowestWorkflowStatus(ctx, &UpdateParentTaskStatusToLowestWorkflowStatus{
				taskRepo:      s.taskRepo,
				Workflows:     project.Workflows,
				ChildrenTasks: childrenTasks,
				ParentTask:    parentTask,
				UpdaterUserID: bsonUserID,
			})
			if serviceErr != nil {
				return nil, serviceErr
			}

			// If Level 1 parent task (Story, Task, Bug) has assignees with points, delete those points from the assignees
			if assigneeesContainPoint(parentTask.Assignees) {
				newAssignees := make([]repositories.UpdateTaskAssigneesRequestAssignee, 0, len(parentTask.Assignees))
				for _, assignee := range parentTask.Assignees {
					newAssignees = append(newAssignees, repositories.UpdateTaskAssigneesRequestAssignee{
						UserID:   assignee.UserID,
						Position: assignee.Position,
					})
				}

				_, err := s.taskRepo.UpdateAssignees(ctx, &repositories.UpdateTaskAssigneesRequest{
					ID:        parentTask.ID,
					Assignees: newAssignees,
					UpdatedBy: bsonUserID,
				})
				if err != nil {
					return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
				}
			}
		}
	}

	return task, nil
}

func (s *taskServiceImpl) GetTaskDetail(ctx context.Context, req *requests.GetTaskDetailPathParam, userId string) (*responses.GetTaskDetailResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
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
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest)
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	task, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.TaskRef, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if task == nil {
		return nil, errutils.NewError(exceptions.ErrTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Task not found: %s", req.TaskRef))
	}

	approvalUserIDs := make([]bson.ObjectID, 0, len(task.Approvals))
	for _, approval := range task.Approvals {
		approvalUserIDs = append(approvalUserIDs, approval.UserID)
	}

	approvals, err := s.userRepo.FindByIDs(ctx, approvalUserIDs)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	approvalMap := make(map[string]models.User, len(approvals))
	for _, approval := range approvals {
		approvalMap[approval.ID.Hex()] = approval
	}

	approvalResponses := make([]responses.GetTaskDetailResponseApprovals, len(task.Approvals))
	for i, approval := range task.Approvals {
		user, ok := approvalMap[approval.UserID.Hex()]
		if !ok {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage("Approval not found")
		}

		var profileUrl = user.DefaultProfileUrl
		if user.UploadedProfileUrl != nil {
			profileUrl = *user.UploadedProfileUrl
		}

		approvalResponses[i] = responses.GetTaskDetailResponseApprovals{
			UserID:      approval.UserID.Hex(),
			Email:       user.Email,
			DisplayName: user.DisplayName,
			ProfileUrl:  profileUrl,
			IsApproved:  approval.IsApproved,
			Reason:      approval.Reason,
		}
	}

	assigneeUserIDs := make([]bson.ObjectID, 0, len(task.Assignees))
	for _, assignee := range task.Assignees {
		if assignee.UserID != nil {
			assigneeUserIDs = append(assigneeUserIDs, *assignee.UserID)
		}
	}

	assignees, err := s.userRepo.FindByIDs(ctx, assigneeUserIDs)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	assigneeMap := make(map[string]models.User, len(assignees))
	for _, assignee := range assignees {
		assigneeMap[assignee.ID.Hex()] = assignee
	}

	assigneeResponses := make([]responses.GetTaskDetailResponseAssignee, len(task.Assignees))
	for i, assignee := range task.Assignees {
		if assignee.UserID == nil {
			assigneeResponses[i] = responses.GetTaskDetailResponseAssignee{
				Position: assignee.Position,
				Point:    assignee.Point,
			}
			continue
		} else if _, ok := assigneeMap[assignee.UserID.Hex()]; !ok {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage("Assignee not found")
		} else {
			user, ok := assigneeMap[assignee.UserID.Hex()]
			if ok {
				var profileUrl = user.DefaultProfileUrl
				if user.UploadedProfileUrl != nil {
					profileUrl = *user.UploadedProfileUrl
				}

				userID := assignee.UserID.Hex()

				assigneeResponses[i] = responses.GetTaskDetailResponseAssignee{
					UserID:      &userID,
					Email:       &user.Email,
					DisplayName: &user.DisplayName,
					ProfileUrl:  &profileUrl,
					Position:    assignee.Position,
					Point:       assignee.Point,
				}
			}
		}
	}

	reporter, err := s.userRepo.FindByID(ctx, task.CreatedBy)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if reporter == nil {
		return nil, errutils.NewError(exceptions.ErrUserNotFound, errutils.InternalServerError).WithDebugMessage(fmt.Sprintf("User not found: %s", task.CreatedBy.Hex()))
	}

	updater, err := s.userRepo.FindByID(ctx, task.UpdatedBy)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if updater == nil {
		return nil, errutils.NewError(exceptions.ErrUserNotFound, errutils.InternalServerError).WithDebugMessage(fmt.Sprintf("User not found: %s", task.UpdatedBy.Hex()))
	}

	var reporterProfileUrl = reporter.DefaultProfileUrl
	if reporter.UploadedProfileUrl != nil {
		reporterProfileUrl = *reporter.UploadedProfileUrl
	}

	var updaterProfileUrl = updater.DefaultProfileUrl
	if updater.UploadedProfileUrl != nil {
		updaterProfileUrl = *updater.UploadedProfileUrl
	}

	return &responses.GetTaskDetailResponse{
		ID:                  task.ID.Hex(),
		TaskRef:             task.TaskRef,
		ProjectID:           task.ProjectID.Hex(),
		Title:               task.Title,
		Description:         task.Description,
		ParentID:            task.ParentID,
		Type:                task.Type,
		Status:              task.Status,
		Priority:            task.Priority,
		Approvals:           approvalResponses,
		Assignees:           assigneeResponses,
		ChildrenPoint:       task.ChildrenPoint,
		HasChildren:         task.HasChildren,
		Sprint:              task.Sprint,
		Attributes:          task.Attributes,
		StartDate:           task.StartDate,
		DueDate:             task.DueDate,
		CreatedAt:           task.CreatedAt,
		ReporterUserID:      task.CreatedBy.Hex(),
		ReporterDisplayName: reporter.DisplayName,
		ReporterProfileUrl:  reporterProfileUrl,
		UpdatedAt:           task.UpdatedAt,
		UpdatedBy:           task.UpdatedBy.Hex(),
		UpdaterDisplayName:  updater.DisplayName,
		UpdaterProfileUrl:   updaterProfileUrl,
	}, nil
}

func (s *taskServiceImpl) GetManyTaskDetail(ctx context.Context, req *requests.GetManyTaskDetailParams, userId string) ([]responses.GetTaskDetailResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
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
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest)
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	tasks, err := s.taskRepo.FindByTaskRefsAndProjectID(ctx, req.TaskRefs, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if len(tasks) == 0 {
		return nil, errutils.NewError(exceptions.ErrTaskNotFound, errutils.BadRequest).WithDebugMessage("Tasks not found")
	}

	response := make([]responses.GetTaskDetailResponse, 0, len(tasks))
	for _, task := range tasks {
		approvalUserIDs := make([]bson.ObjectID, 0, len(task.Approvals))
		for _, approval := range task.Approvals {
			approvalUserIDs = append(approvalUserIDs, approval.UserID)
		}

		approvals, err := s.userRepo.FindByIDs(ctx, approvalUserIDs)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		approvalMap := make(map[string]models.User, len(approvals))
		for _, approval := range approvals {
			approvalMap[approval.ID.Hex()] = approval
		}

		approvalResponses := make([]responses.GetTaskDetailResponseApprovals, len(task.Approvals))
		for i, approval := range task.Approvals {
			user, ok := approvalMap[approval.UserID.Hex()]
			if !ok {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage("Approval not found")
			}

			var profileUrl = user.DefaultProfileUrl
			if user.UploadedProfileUrl != nil {
				profileUrl = *user.UploadedProfileUrl
			}

			approvalResponses[i] = responses.GetTaskDetailResponseApprovals{
				UserID:      approval.UserID.Hex(),
				Email:       user.Email,
				DisplayName: user.DisplayName,
				ProfileUrl:  profileUrl,
				IsApproved:  approval.IsApproved,
				Reason:      approval.Reason,
			}
		}

		assigneeUserIDs := make([]bson.ObjectID, 0, len(task.Assignees))
		for _, assignee := range task.Assignees {
			if assignee.UserID != nil {
				assigneeUserIDs = append(assigneeUserIDs, *assignee.UserID)
			}
		}

		assignees, err := s.userRepo.FindByIDs(ctx, assigneeUserIDs)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		assigneeMap := make(map[string]models.User, len(assignees))
		for _, assignee := range assignees {
			assigneeMap[assignee.ID.Hex()] = assignee
		}

		assigneeResponses := make([]responses.GetTaskDetailResponseAssignee, len(task.Assignees))
		for i, assignee := range task.Assignees {
			if assignee.UserID == nil {
				assigneeResponses[i] = responses.GetTaskDetailResponseAssignee{
					Position: assignee.Position,
					Point:    assignee.Point,
				}
				continue
			} else if _, ok := assigneeMap[assignee.UserID.Hex()]; !ok {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage("Assignee not found")
			} else {
				user, ok := assigneeMap[assignee.UserID.Hex()]
				if ok {
					var profileUrl = user.DefaultProfileUrl
					if user.UploadedProfileUrl != nil {
						profileUrl = *user.UploadedProfileUrl
					}

					userID := assignee.UserID.Hex()

					assigneeResponses[i] = responses.GetTaskDetailResponseAssignee{
						UserID:      &userID,
						Email:       &user.Email,
						DisplayName: &user.DisplayName,
						ProfileUrl:  &profileUrl,
						Position:    assignee.Position,
						Point:       assignee.Point,
					}
				}
			}
		}

		reporter, err := s.userRepo.FindByID(ctx, task.CreatedBy)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if reporter == nil {
			return nil, errutils.NewError(exceptions.ErrUserNotFound, errutils.InternalServerError).WithDebugMessage(fmt.Sprintf("User not found: %s", task.CreatedBy.Hex()))
		}

		updater, err := s.userRepo.FindByID(ctx, task.UpdatedBy)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if updater == nil {
			return nil, errutils.NewError(exceptions.ErrUserNotFound, errutils.InternalServerError).WithDebugMessage(fmt.Sprintf("User not found: %s", task.UpdatedBy.Hex()))
		}

		var reporterProfileUrl = reporter.DefaultProfileUrl
		if reporter.UploadedProfileUrl != nil {
			reporterProfileUrl = *reporter.UploadedProfileUrl
		}

		var updaterProfileUrl = updater.DefaultProfileUrl
		if updater.UploadedProfileUrl != nil {
			updaterProfileUrl = *updater.UploadedProfileUrl
		}

		response = append(response, responses.GetTaskDetailResponse{
			ID:                  task.ID.Hex(),
			TaskRef:             task.TaskRef,
			ProjectID:           task.ProjectID.Hex(),
			Title:               task.Title,
			Description:         task.Description,
			ParentID:            task.ParentID,
			Type:                task.Type,
			Status:              task.Status,
			Priority:            task.Priority,
			Approvals:           approvalResponses,
			Assignees:           assigneeResponses,
			ChildrenPoint:       task.ChildrenPoint,
			HasChildren:         task.HasChildren,
			Sprint:              task.Sprint,
			Attributes:          task.Attributes,
			StartDate:           task.StartDate,
			DueDate:             task.DueDate,
			CreatedAt:           task.CreatedAt,
			ReporterUserID:      task.CreatedBy.Hex(),
			ReporterDisplayName: reporter.DisplayName,
			ReporterProfileUrl:  reporterProfileUrl,
			UpdatedAt:           task.UpdatedAt,
			UpdatedBy:           task.UpdatedBy.Hex(),
			UpdaterDisplayName:  updater.DisplayName,
			UpdaterProfileUrl:   updaterProfileUrl,
		})
	}

	return response, nil
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

	var taskTypes []models.TaskType
	for _, taskType := range req.Types {
		if !models.TaskType(taskType).IsValid() {
			return nil, errutils.NewError(exceptions.ErrInvalidTaskType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid task type: %s", taskType))
		}
		taskTypes = append(taskTypes, models.TaskType(taskType))
	}

	var (
		isTaskWithNoSprint bool
		bsonSprintIDs      []bson.ObjectID
	)
	if req.IsTaskInBacklog != nil && *req.IsTaskInBacklog {
		isTaskWithNoSprint = true
	} else if req.SprintIDs != nil {
		for _, sprintID := range req.SprintIDs {
			bsonSprintID, err := bson.ObjectIDFromHex(sprintID)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
			}
			bsonSprintIDs = append(bsonSprintIDs, bsonSprintID)
		}
	}

	var (
		bsonParentID     *bson.ObjectID
		isTaskWithNoEpic bool
	)
	if req.EpicTaskID != nil {
		if *req.EpicTaskID == constant.SearchTaskParamsTaskWithNoEpicFilter {
			isTaskWithNoEpic = true
		} else {
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
		ProjectID:          bsonProjectID,
		TaskTypes:          taskTypes,
		SprintIDs:          bsonSprintIDs,
		IsTaskWithNoSprint: isTaskWithNoSprint,
		EpicTaskID:         bsonParentID,
		IsTaskWithNoEpic:   isTaskWithNoEpic,
		UserIDs:            userIDs,
		Positions:          req.Positions,
		Statuses:           req.Statuses,
		IsDoneStatuses:     getDoneStatuses(project),
		SearchKeyword:      req.SearchKeyword,
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
		assigneesUserIDs := make([]bson.ObjectID, 0, len(task.Assignees))
		for _, assignee := range task.Assignees {
			if assignee.UserID != nil {
				assigneesUserIDs = append(assigneesUserIDs, *assignee.UserID)
			}
		}

		assignees, err := s.userRepo.FindByIDs(ctx, assigneesUserIDs)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		assigneeMap := make(map[string]models.User, len(assignees))
		for _, assignee := range assignees {
			assigneeMap[assignee.ID.Hex()] = assignee
		}

		assigneeResponses := make([]responses.SearchTaskResponseAssignee, len(task.Assignees))
		for i, assignee := range task.Assignees {
			if assignee.UserID == nil {
				assigneeResponses[i] = responses.SearchTaskResponseAssignee{
					Position: assignee.Position,
					Point:    assignee.Point,
				}
				continue
			} else if _, ok := assigneeMap[assignee.UserID.Hex()]; !ok {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage("Assignee not found")
			} else {
				user, ok := assigneeMap[assignee.UserID.Hex()]
				if ok {
					var profileUrl = user.DefaultProfileUrl
					if user.UploadedProfileUrl != nil {
						profileUrl = *user.UploadedProfileUrl
					}

					userID := assignee.UserID.Hex()

					assigneeResponses[i] = responses.SearchTaskResponseAssignee{
						UserID:      &userID,
						Email:       &user.Email,
						DisplayName: &user.DisplayName,
						ProfileUrl:  &profileUrl,
						Position:    assignee.Position,
						Point:       assignee.Point,
					}
				}
			}
		}

		var parentTitleResp *string
		if task.ParentID != nil {
			parentTitle, ok := parentTasksMap[task.ParentID.Hex()]
			if !ok {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage("Parent task not found")
			}

			parentTitleResp = parentTitle
		}

		response = append(response, responses.SearchTaskResponse{
			ID:            task.ID.Hex(),
			TaskRef:       task.TaskRef,
			Title:         task.Title,
			ParentID:      conv.BsonObjectIDPtrToStringPtr(task.ParentID),
			ParentTitle:   parentTitleResp,
			Type:          task.Type.String(),
			Status:        task.Status,
			Priority:      task.Priority.String(),
			Assignees:     assigneeResponses,
			Approvals:     task.Approvals,
			ChildrenPoint: task.ChildrenPoint,
			HasChildren:   task.HasChildren,
			Sprint:        task.Sprint,
		})
	}

	return response, nil
}

func (s *taskServiceImpl) GetChildrenTasks(ctx context.Context, req *requests.GetChildrenTasksParams, userId string) ([]responses.GetChildrenTasksResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
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
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest)
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	parentTask, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.ParentTaskRef, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	tasks, err := s.taskRepo.FindByParentID(ctx, parentTask.ID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if len(tasks) == 0 {
		return []responses.GetChildrenTasksResponse{}, nil
	}

	response := make([]responses.GetChildrenTasksResponse, 0, len(tasks))
	for _, task := range tasks {
		approvalUserIDs := make([]bson.ObjectID, 0, len(task.Approvals))
		for _, approval := range task.Approvals {
			approvalUserIDs = append(approvalUserIDs, approval.UserID)
		}

		approvals, err := s.userRepo.FindByIDs(ctx, approvalUserIDs)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		approvalMap := make(map[string]models.User, len(approvals))
		for _, approval := range approvals {
			approvalMap[approval.ID.Hex()] = approval
		}

		approvalResponses := make([]responses.GetTaskDetailResponseApprovals, len(task.Approvals))
		for i, approval := range task.Approvals {
			user, ok := approvalMap[approval.UserID.Hex()]
			if !ok {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage("Approval not found")
			}

			var profileUrl = user.DefaultProfileUrl
			if user.UploadedProfileUrl != nil {
				profileUrl = *user.UploadedProfileUrl
			}

			approvalResponses[i] = responses.GetTaskDetailResponseApprovals{
				UserID:      approval.UserID.Hex(),
				Email:       user.Email,
				DisplayName: user.DisplayName,
				ProfileUrl:  profileUrl,
				IsApproved:  approval.IsApproved,
				Reason:      approval.Reason,
			}
		}

		assigneeUserIDs := make([]bson.ObjectID, 0, len(task.Assignees))
		for _, assignee := range task.Assignees {
			if assignee.UserID != nil {
				assigneeUserIDs = append(assigneeUserIDs, *assignee.UserID)
			}
		}

		assignees, err := s.userRepo.FindByIDs(ctx, assigneeUserIDs)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}

		assigneeMap := make(map[string]models.User, len(assignees))
		for _, assignee := range assignees {
			assigneeMap[assignee.ID.Hex()] = assignee
		}

		assigneeResponses := make([]responses.GetTaskDetailResponseAssignee, len(task.Assignees))
		for i, assignee := range task.Assignees {
			if assignee.UserID == nil {
				assigneeResponses[i] = responses.GetTaskDetailResponseAssignee{
					Position: assignee.Position,
					Point:    assignee.Point,
				}
				continue
			} else if _, ok := assigneeMap[assignee.UserID.Hex()]; !ok {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage("Assignee not found")
			} else {
				user, ok := assigneeMap[assignee.UserID.Hex()]
				if ok {
					var profileUrl = user.DefaultProfileUrl
					if user.UploadedProfileUrl != nil {
						profileUrl = *user.UploadedProfileUrl
					}

					userID := assignee.UserID.Hex()

					assigneeResponses[i] = responses.GetTaskDetailResponseAssignee{
						UserID:      &userID,
						Email:       &user.Email,
						DisplayName: &user.DisplayName,
						ProfileUrl:  &profileUrl,
						Position:    assignee.Position,
						Point:       assignee.Point,
					}
				}
			}
		}

		response = append(response, responses.GetChildrenTasksResponse{
			ID:        task.ID.Hex(),
			TaskRef:   task.TaskRef,
			ProjectID: task.ProjectID.Hex(),
			Title:     task.Title,
			Type:      task.Type.String(),
			Status:    task.Status,
			Priority:  task.Priority.String(),
			Approvals: approvalResponses,
			Assignees: assigneeResponses,
		})
	}

	return response, nil
}

func (s *taskServiceImpl) UpdateDetail(ctx context.Context, req *requests.UpdateTaskDetailRequest, userId string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	// Check if task priority is valid
	if !models.TaskPriority(req.Priority).IsValid() {
		return nil, errutils.NewError(exceptions.ErrInvalidTaskPriority, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid task priority: %s", req.Priority))
	}

	task, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.TaskRef, bsonProjectID)
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

	updatedTask, err := s.taskRepo.UpdateDetail(ctx, &repositories.UpdateTaskDetailRequest{
		ID:          task.ID,
		Title:       req.Title,
		Description: req.Description,
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

func (s *taskServiceImpl) UpdateTitle(ctx context.Context, req *requests.UpdateTaskTitleRequest, userId string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.TaskRef, bsonProjectID)
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

	updatedTask, err := s.taskRepo.UpdateTitle(ctx, &repositories.UpdateTaskTitleRequest{
		ID:        task.ID,
		Title:     req.Title,
		UpdatedBy: bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return updatedTask, nil
}

func (s *taskServiceImpl) UpdateParentID(ctx context.Context, req *requests.UpdateTaskParentIdRequest, userID string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.TaskRef, bsonProjectID)
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
		isParentTaskChanged      = (task.ParentID == nil && req.ParentID != nil) || (req.ParentID != nil && task.ParentID != nil && *req.ParentID != task.ParentID.Hex())
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
			// Update Task's Sprint to Parent Task's Sprint
			_, err := s.taskRepo.UpdateCurrentSprintID(ctx, &repositories.UpdateTaskCurrentSprintIDRequest{
				ID:              task.ID,
				CurrentSprintID: newParentTask.Sprint.CurrentSprintID,
				UpdatedBy:       bsonUserID,
			})
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}

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

			if newParentTask.ParentID != nil {
				epicParentTask, err := s.taskRepo.FindByID(ctx, *newParentTask.ParentID)
				if err != nil {
					return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
				} else if epicParentTask == nil {
					return nil, errutils.NewError(exceptions.ErrParentTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task not found: %s", *newParentTask.ParentID))
				}

				_, err = s.taskRepo.UpdateChildrenPoint(ctx, &repositories.UpdateTaskChildrenPointRequest{
					ID:            *newParentTask.ParentID,
					ChildrenPoint: epicParentTask.ChildrenPoint + totalPoint,
				})
				if err != nil {
					return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
				}
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

	updatedTask, err := s.taskRepo.UpdateParentID(ctx, &repositories.UpdateTaskParentIDRequest{
		ID:        task.ID,
		ParentID:  nullableBsonTaskParentID,
		UpdatedBy: bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return updatedTask, nil
}

func (s *taskServiceImpl) UpdateType(ctx context.Context, req *requests.UpdateTaskTypeRequest, userId string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	if !models.TaskType(req.Type).IsValid() {
		return nil, errutils.NewError(exceptions.ErrInvalidTaskType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Invalid task type: %s", req.Type))
	}

	task, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.TaskRef, bsonProjectID)
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

	/*
		Currently, only task in the same level can be converted to each other.

		(Task <=> Story <=> Bug)
	*/

	if task.Type == models.TaskTypeEpic || task.Type == models.TaskTypeSubTask ||
		req.Type == models.TaskTypeEpic.String() || req.Type == models.TaskTypeSubTask.String() {
		return nil, errutils.NewError(exceptions.ErrOnlyTaskInTheSameLevelCanChangeType, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Task type: %s", task.Type))
	}

	updatedTask, err := s.taskRepo.UpdateType(ctx, &repositories.UpdateTaskTypeRequest{
		ID:        task.ID,
		Type:      models.TaskType(req.Type),
		UpdatedBy: bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return updatedTask, nil
}

func (s *taskServiceImpl) UpdateStatus(ctx context.Context, req *requests.UpdateTaskStatusRequest, userId string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
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
		return nil, errutils.NewError(exceptions.ErrProjectNotFound, errutils.BadRequest)
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	task, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.TaskID, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if task == nil {
		return nil, errutils.NewError(exceptions.ErrTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Task not found: %s", req.TaskID))
	}

	validStatuses := make([]string, 0, len(project.Workflows))
	for _, workflow := range project.Workflows {
		validStatuses = append(validStatuses, workflow.Status)
	}

	if !array.ContainAny(validStatuses, []string{req.Status}) {
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

	// If the updated task is the level 1 task (Task, Story, Bug)
	// Update all children tasks' status to the updated parent task's status
	if array.ContainAny(
		[]string{task.Type.String()},
		[]string{models.TaskTypeStory.String(), models.TaskTypeTask.String(), models.TaskTypeBug.String()},
	) {
		childrenTasks, err := s.taskRepo.FindByParentID(ctx, task.ID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if len(childrenTasks) != 0 {
			childrenTaskBsonIDs := make([]bson.ObjectID, 0, len(childrenTasks))
			for _, childrenTask := range childrenTasks {
				childrenTaskBsonIDs = append(childrenTaskBsonIDs, childrenTask.ID)
			}

			err = s.taskRepo.UpdateManyTasksStatus(ctx, &repositories.UpdateManyTasksStatusRequest{
				TaskIDs:   childrenTaskBsonIDs,
				Status:    req.Status,
				UpdatedBy: bsonUserID,
			})
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}
		}

	} else if updatedTask.Type == models.TaskTypeSubTask {
		// If the updated task is a SubTask
		// Update the parent task's status to the lowest status of all children tasks
		if updatedTask.ParentID != nil {
			parentTask, err := s.taskRepo.FindByID(ctx, *updatedTask.ParentID)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			} else if parentTask == nil {
				return nil, errutils.NewError(exceptions.ErrParentTaskNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Parent task not found: %s", *updatedTask.ParentID))
			}

			childrenTasks, err := s.taskRepo.FindByParentID(ctx, parentTask.ID)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}

			serviceErr := updateParentTaskStatusToLowestWorkflowStatus(ctx, &UpdateParentTaskStatusToLowestWorkflowStatus{
				taskRepo:      s.taskRepo,
				Workflows:     project.Workflows,
				ChildrenTasks: childrenTasks,
				ParentTask:    parentTask,
				UpdaterUserID: bsonUserID,
			})
			if serviceErr != nil {
				return nil, serviceErr
			}
		}
	}

	return updatedTask, nil
}

func (s *taskServiceImpl) UpdateApprovals(ctx context.Context, req *requests.UpdateTaskApprovalsRequest, userID string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.TaskRef, bsonProjectID)
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

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.TaskRef, bsonProjectID)
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

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.TaskRef, bsonProjectID)
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

	assignees := make([]repositories.UpdateTaskAssigneesRequestAssignee, 0, len(req.Assignees))
	for _, assignee := range req.Assignees {
		var bsonAssigneeUserID *bson.ObjectID
		if assignee.UserId != nil {
			assigneeUserID, err := bson.ObjectIDFromHex(*assignee.UserId)
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
			}
			bsonAssigneeUserID = &assigneeUserID
		}

		assigneeRequest := repositories.UpdateTaskAssigneesRequestAssignee{
			Position: assignee.Position,
			UserID:   bsonAssigneeUserID,
			Point:    assignee.Point,
		}

		assignees = append(assignees, assigneeRequest)
	}

	if task.Type == models.TaskTypeSubTask ||
		array.ContainAny(
			[]string{task.Type.String()},
			[]string{
				models.TaskTypeStory.String(),
				models.TaskTypeTask.String(),
				models.TaskTypeBug.String(),
			},
		) {
		var currentTotalPoint int
		for _, assignee := range task.Assignees {
			if assignee.Point != nil {
				currentTotalPoint += *assignee.Point
			}
		}

		serviceErr := updateParentTaskPoints(ctx, s.taskRepo, task, assignees, currentTotalPoint)
		if serviceErr != nil {
			return nil, serviceErr
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

func (s *taskServiceImpl) UpdateSprint(ctx context.Context, req *requests.UpdateTaskSprintRequest, userId string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.TaskRef, bsonProjectID)
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
		bsonCurrentSprintID *bson.ObjectID
		startDate           *time.Time
		endDate             *time.Time
	)
	if req.CurrentSprintID != nil {
		currentSprintID, err := bson.ObjectIDFromHex(*req.CurrentSprintID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
		}
		bsonCurrentSprintID = &currentSprintID

		sprint, err := s.sprintRepo.FindByID(ctx, *bsonCurrentSprintID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		} else if sprint == nil {
			return nil, errutils.NewError(exceptions.ErrSprintNotFound, errutils.BadRequest).WithDebugMessage(fmt.Sprintf("Sprint not found: %s", *req.CurrentSprintID))
		}

		if sprint.StartDate != nil && sprint.EndDate != nil {
			startDate = sprint.StartDate
			endDate = sprint.EndDate
		}
	}

	updatedTask, err := s.taskRepo.UpdateCurrentSprintID(ctx, &repositories.UpdateTaskCurrentSprintIDRequest{
		ID:              task.ID,
		CurrentSprintID: bsonCurrentSprintID,
		UpdatedBy:       bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	if startDate != nil && endDate != nil {
		_, err = s.taskRepo.UpdateStartDateAndDueDate(ctx, &repositories.UpdateTaskStartDateAndDueDateRequest{
			ID:        task.ID,
			StartDate: startDate,
			DueDate:   endDate,
			UpdatedBy: bsonUserID,
		})
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}
	}

	// If the updated task is the level 1 task (Task, Story, Bug)
	// Update all children tasks' sprint to the updated parent task's sprint
	if array.ContainAny(
		[]string{task.Type.String()},
		[]string{
			models.TaskTypeStory.String(),
			models.TaskTypeTask.String(),
			models.TaskTypeBug.String(),
		},
	) {
		childrenTasks, err := s.taskRepo.FindByParentID(ctx, task.ID)
		if err != nil {
			return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
		}
		if len(childrenTasks) > 0 {
			childrenTaskIDs := make([]bson.ObjectID, 0, len(childrenTasks))
			for _, childrenTask := range childrenTasks {
				childrenTaskIDs = append(childrenTaskIDs, childrenTask.ID)
			}

			err = s.taskRepo.BulkUpdateCurrentSprintID(ctx, &repositories.BulkUpdateCurrentSprintIDRequest{
				TaskIDs:         childrenTaskIDs,
				CurrentSprintID: bsonCurrentSprintID,
				UpdatedBy:       bsonUserID,
			})
			if err != nil {
				return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
			}

			// Update all children tasks' start date and due date to the updated parent task's sprint start date and end date
			if startDate != nil && endDate != nil {
				err = s.taskRepo.BulkUpdateStartDateAndDueDate(ctx, &repositories.BulkUpdateStartDateAndDueDateRequest{
					TaskIDs:   childrenTaskIDs,
					StartDate: startDate,
					DueDate:   endDate,
					UpdatedBy: bsonUserID,
				})
				if err != nil {
					return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
				}
			}
		}
	}

	return updatedTask, nil
}

func (s *taskServiceImpl) UpdateAttributes(ctx context.Context, req *requests.UpdateTaskAttributesRequest, userID string) (*models.Task, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.TaskRef, bsonProjectID)
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
			value, err = strconv.ParseFloat(attribute.Value, 64)
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

// To be further implemented (prompt)
func (s *taskServiceImpl) GenerateDescription(ctx context.Context, req *requests.GenerateDescriptionRequest, userId string) (*responses.GenerateDescriptionResponse, *errutils.Error) {

	prompt := fmt.Sprintf(`
	generate task description base on the task title "%s"  in this structure

   [

           {

             id: "4bbfc57b-d00c-49f6-af09-caf198534f1f",

             type: "paragraph",

             props: { textColor: "default", backgroundColor: "default", textAlignment: "left" },

             content: [{ type: "text", text: " ", styles: {} }],

             children: [],

           },

          ....

         ]


		Structure Reference

		type ParagraphBlock = {  id: string;  type: "paragraph";  props: DefaultProps;  content: InlineContent[];  children: Block[];};



		type HeadingBlock = {  id: string;  type: "heading";  props: {    level: 1 | 2 | 3 = 1;  } & DefaultProps;  content: InlineContent[];  children: Block[];};



		type BulletListItemBlock = {  id: string;  type: "bulletListItem";  props: DefaultProps;  content: InlineContent[];  children: Block[];};



		type NumberedListItemBlock = {  id: string;  type: "numberedListItem";  props: DefaultProps;  content: InlineContent[];  children: Block[];};

		response only JSON
	`, req.Prompt)

	resp, err := s.geminiRepo.GenerateTaskDescription(ctx, prompt)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	var parts []genai.Part
	for _, content := range resp {
		parts = append(parts, content.Parts...)
	}

	return &responses.GenerateDescriptionResponse{
		Description: parts,
	}, nil
}
