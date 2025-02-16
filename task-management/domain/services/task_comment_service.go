package services

import (
	"context"

	"github.com/cnc-csku/task-nexus/go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskCommentService interface {
	Create(ctx context.Context, req *requests.CreateTaskCommentRequest, userID string) (*models.TaskComment, *errutils.Error)
}

type taskCommentServiceImpl struct {
	taskCommentRepo   repositories.TaskCommentRepository
	taskRepo          repositories.TaskRepository
	projectRepo       repositories.ProjectRepository
	projectMemberRepo repositories.ProjectMemberRepository
}

func NewTaskCommentService(
	taskCommentRepo repositories.TaskCommentRepository,
	taskRepo repositories.TaskRepository,
	projectRepo repositories.ProjectRepository,
	projectMemberRepo repositories.ProjectMemberRepository,
) TaskCommentService {
	return &taskCommentServiceImpl{
		taskCommentRepo:   taskCommentRepo,
		taskRepo:          taskRepo,
		projectRepo:       projectRepo,
		projectMemberRepo: projectMemberRepo,
	}
}

func (s *taskCommentServiceImpl) Create(ctx context.Context, req *requests.CreateTaskCommentRequest, userID string) (*models.TaskComment, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskID(ctx, req.TaskID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if task == nil {
		return nil, errutils.NewError(exceptions.ErrTaskNotFound, errutils.BadRequest).WithDebugMessage("task not found")
	}

	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, task.ProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	comment, err := s.taskCommentRepo.Create(ctx, &repositories.CreateTaskCommentRequest{
		TaskID:  req.TaskID,
		Content: req.Content,
		UserID:  bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return comment, nil
}
