package grpcserver

import (
	"context"

	taskv1 "github.com/cnc-csku/task-nexus/api-specification/gen/proto/task/v1"
	"github.com/cnc-csku/task-nexus/task-management/domain/services"
)

type MemberHandler struct {
	taskv1.UnimplementedMemberServiceServer
	services.MemberUseCase
}

func (h *MemberHandler) GetMembers(ctx context.Context, in *taskv1.GetMembersRequest) (*taskv1.GetMembersResponse, error) {
	resp, err := h.MemberUseCase.GetMembers(ctx, in)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
