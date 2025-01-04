package services

import (
	"context"

	"github.com/cnc-csku/task-nexus/task-management/domain/repositories"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/repositories/grpcclient"
)

type MemberService interface {
	GetMembers(ctx context.Context, in *requests.GetMembersRequest) (*responses.GetMembersResponse, error)
}

type memberService struct {
	memberRepo repositories.MemberRepository
	grpcClient *grpcclient.GrpcClient
}

func NewMemberService(
	memberRepo repositories.MemberRepository,
	grpcClient *grpcclient.GrpcClient,
) MemberService {
	return &memberService{
		memberRepo: memberRepo,
		grpcClient: grpcClient,
	}
}

func (u *memberService) GetMembers(ctx context.Context, in *requests.GetMembersRequest) (*responses.GetMembersResponse, error) {
	return nil, nil
}
