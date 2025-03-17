package services

import (
	"context"

	"github.com/cnc-csku/task-nexus-go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskCommentService interface {
	Create(ctx context.Context, req *requests.CreateTaskCommentRequest, userID string) (*models.TaskComment, *errutils.Error)
	List(ctx context.Context, req *requests.ListTaskCommentPathParams, userID string) ([]responses.ListTaskCommentResponse, *errutils.Error)
}

type taskCommentServiceImpl struct {
	userRepo          repositories.UserRepository
	taskCommentRepo   repositories.TaskCommentRepository
	taskRepo          repositories.TaskRepository
	projectRepo       repositories.ProjectRepository
	projectMemberRepo repositories.ProjectMemberRepository
}

func NewTaskCommentService(
	userRepo repositories.UserRepository,
	taskCommentRepo repositories.TaskCommentRepository,
	taskRepo repositories.TaskRepository,
	projectRepo repositories.ProjectRepository,
	projectMemberRepo repositories.ProjectMemberRepository,
) TaskCommentService {
	return &taskCommentServiceImpl{
		userRepo:          userRepo,
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

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.TaskRef, bsonProjectID)
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
		TaskID:  task.ID,
		Content: req.Content,
		UserID:  bsonUserID,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return comment, nil
}

func (s *taskCommentServiceImpl) List(ctx context.Context, req *requests.ListTaskCommentPathParams, userID string) ([]responses.ListTaskCommentResponse, *errutils.Error) {
	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	task, err := s.taskRepo.FindByTaskRefAndProjectID(ctx, req.TaskRef, bsonProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if task == nil {
		return nil, errutils.NewError(exceptions.ErrTaskNotFound, errutils.BadRequest).WithDebugMessage("task not found")
	}

	// Check if the user is a member of the project
	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, task.ProjectID, bsonUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.BadRequest).WithDebugMessage("User is not a member of the project")
	}

	comments, err := s.taskCommentRepo.FindByTaskID(ctx, task.ID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	users, err := s.userRepo.FindByIDs(ctx, extractUserIDsFromComments(comments))
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	}

	return buildTaskComments(comments, mapUsersByID(users)), nil
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

func buildTaskComments(comments []*models.TaskComment, userMap map[string]string) []responses.ListTaskCommentResponse {
	taskComments := make([]responses.ListTaskCommentResponse, 0, len(comments))
	for _, comment := range comments {
		taskComments = append(taskComments, responses.ListTaskCommentResponse{
			ID:              comment.ID.Hex(),
			Content:         comment.Content,
			UserID:          comment.UserID.Hex(),
			UserDisplayName: userMap[comment.UserID.Hex()],
			UserProfileUrl:  userMap[comment.UserID.Hex()],
			TaskID:          comment.TaskID.Hex(),
			CreatedAt:       comment.CreatedAt,
			UpdatedAt:       comment.UpdatedAt,
		})
	}
	return taskComments
}
