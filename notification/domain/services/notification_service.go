package services

import (
	"context"
	"fmt"

	notificationv1 "github.com/cnc-csku/task-nexus/api-specification/gen/proto/notification/v1"
	core_grpcclient "github.com/cnc-csku/task-nexus/go-lib/grpcclient"
)

type NotificationService interface {
	TestConnection(ctx context.Context, in *notificationv1.TestConnectionRequest) (*notificationv1.TestConnectionResponse, error)
}

type notificationService struct {
	grpcClient *core_grpcclient.GrpcClient
}

func NewNotificationService(
	grpcClient *core_grpcclient.GrpcClient,
) NotificationService {
	return &notificationService{
		grpcClient: grpcClient,
	}
}

func (u *notificationService) TestConnection(ctx context.Context, in *notificationv1.TestConnectionRequest) (*notificationv1.TestConnectionResponse, error) {
	fmt.Println("Notification 1")
	return &notificationv1.TestConnectionResponse{
		Message: "Success",
	}, nil
}
