package services

import (
	"context"

	"github.com/cnc-csku/task-nexus-go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/domain/exceptions"
	"github.com/cnc-csku/task-nexus/task-management/domain/models"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProjectMemberService interface {
	UpdatePosition(ctx context.Context, req *requests.UpdateMemberPositionRequest, userID string) (*models.ProjectMember, *errutils.Error)
}

type projectMemberServiceImpl struct {
	userRepo          repositories.UserRepository
	projectRepo       repositories.ProjectRepository
	projectMemberRepo repositories.ProjectMemberRepository
}

func NewProjectMemberService(
	userRepo repositories.UserRepository,
	projectRepo repositories.ProjectRepository,
	projectMemberRepo repositories.ProjectMemberRepository,
) ProjectMemberService {
	return &projectMemberServiceImpl{
		userRepo:          userRepo,
		projectRepo:       projectRepo,
		projectMemberRepo: projectMemberRepo,
	}
}

func (s *projectMemberServiceImpl) UpdatePosition(ctx context.Context, req *requests.UpdateMemberPositionRequest, userID string) (*models.ProjectMember, *errutils.Error) {
	bsonRequesterUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}
	bsonProjectID, err := bson.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}
	bsonTargetUserID, err := bson.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.BadRequest).WithDebugMessage(err.Error())
	}

	// Check if the requester is the owner or moderator of the project
	member, err := s.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonRequesterUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrMemberNotFoundInProject, errutils.BadRequest).WithDebugMessage("Member not found in project")
	} else if member.Role != models.ProjectMemberRoleOwner && member.Role != models.ProjectMemberRoleModerator {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.Forbidden).WithDebugMessage("Only owner and moderator can update member position")
	} else if req.Position == string(models.ProjectMemberRoleOwner) && member.Role != models.ProjectMemberRoleOwner {
		return nil, errutils.NewError(exceptions.ErrPermissionDenied, errutils.Forbidden).WithDebugMessage("Only owner can update member position to owner")
	}

	// Check if the member exists
	member, err = s.projectMemberRepo.FindByProjectIDAndUserID(ctx, bsonProjectID, bsonTargetUserID)
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrMemberNotFoundInProject, errutils.BadRequest).WithDebugMessage("Member not found in project")
	}

	// Update the member position
	member, err = s.projectMemberRepo.UpdatePositionByID(ctx, &repositories.UpdatePositionRequest{
		ID:       member.ID,
		Position: req.Position,
	})
	if err != nil {
		return nil, errutils.NewError(exceptions.ErrInternalError, errutils.InternalServerError).WithDebugMessage(err.Error())
	} else if member == nil {
		return nil, errutils.NewError(exceptions.ErrMemberNotFoundInProject, errutils.BadRequest).WithDebugMessage("Member not found in project")
	}

	return member, nil
}
