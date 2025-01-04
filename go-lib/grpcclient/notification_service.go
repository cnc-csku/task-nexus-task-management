package core_grpcclient

import (
	"context"
	"fmt"
	"log"

	notificationv1 "github.com/cnc-csku/task-nexus/api-specification/gen/proto/notification/v1"
)

func (g *GrpcClient) WithNotificationServiceClient(ctx context.Context) {
	log.Println("✅ Initializing NotificationService gRPC client")

	if g.config.NotificationService.Host == "" {
		log.Println("NotificationService is not configured")
		return
	}

	conn, err := newGrpcConnection(g.config.NotificationService)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to NotificationService: %v", err))
	}

	grpcClientList[g.config.NotificationService.Name] = conn

	g.NotificationService = notificationv1.NewNotificationServiceClient(conn)

	log.Println("✅ NotificationService gRPC client initialized")
}
