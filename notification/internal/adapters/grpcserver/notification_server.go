package grpcserver

import (
	"context"

	notificationv1 "github.com/cnc-csku/task-nexus/api-specification/gen/proto/notification/v1"
	"github.com/cnc-csku/task-nexus/notification/domain/services"
)

type NotificationServer struct {
	notificationv1.UnimplementedNotificationServiceServer
	services.NotificationService
}

func (s *NotificationServer) TestConnection(ctx context.Context, in *notificationv1.TestConnectionRequest) (*notificationv1.TestConnectionResponse, error) {
	resp, err := s.NotificationService.TestConnection(ctx, in)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
