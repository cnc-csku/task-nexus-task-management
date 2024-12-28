package services

import (
	"context"

	taskv1 "github.com/cnc-csku/task-nexus/api-specification/gen/proto/task/v1"
	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
)

type MemberUseCase interface {
	GetMembers(ctx context.Context, in *taskv1.GetMembersRequest) (*taskv1.GetMembersResponse, error)
}

type memberUseCase struct {
	memberRepo repositories.MemberRepository
}

func NewMemberUseCase(memberRepo repositories.MemberRepository) MemberUseCase {
	return &memberUseCase{memberRepo}
}

func (u *memberUseCase) GetMembers(ctx context.Context, in *taskv1.GetMembersRequest) (*taskv1.GetMembersResponse, error) {
	return nil, nil
}
