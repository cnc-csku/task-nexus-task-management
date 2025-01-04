package services

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
)

type MemberService interface {
	GetMembers(ctx context.Context, in *requests.GetMembersRequest) (*responses.GetMembersResponse, error)
}

type memberService struct {
	memberRepo repositories.MemberRepository
}

func NewMemberService(memberRepo repositories.MemberRepository) MemberService {
	return &memberService{memberRepo}
}

func (u *memberService) GetMembers(ctx context.Context, in *requests.GetMembersRequest) (*responses.GetMembersResponse, error) {
	return nil, nil
}
