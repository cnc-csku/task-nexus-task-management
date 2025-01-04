package services

import (
	"context"

	notificationv1 "github.com/cnc-csku/task-nexus/api-specification/gen/proto/notification/v1"
	"github.com/cnc-csku/task-nexus/task-management/domain/requests"
	"github.com/cnc-csku/task-nexus/task-management/domain/responses"
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/repositories/grpcclient"
)

type CommonService interface {
	TestNotification(ctx context.Context, in *requests.TestNotificationRequest) (*responses.TestNotificationResponse, error)
}

type commonService struct {
	grpcClient *grpcclient.GrpcClient
}

func NewCommonService(
	grpcClient *grpcclient.GrpcClient,
) CommonService {
	return &commonService{
		grpcClient: grpcClient,
	}
}

func (u *commonService) TestNotification(ctx context.Context, in *requests.TestNotificationRequest) (*responses.TestNotificationResponse, error) {
	resp, err := u.grpcClient.Grpcclient.NotificationService.TestConnection(ctx, &notificationv1.TestConnectionRequest{})
	if err != nil {
		return nil, err
	}

	return &responses.TestNotificationResponse{
		Message: resp.Message,
	}, nil
}
