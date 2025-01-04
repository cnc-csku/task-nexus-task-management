package core_grpcclient

import (
	"context"
	"fmt"
	"log"

	taskv1 "github.com/cnc-csku/task-nexus/api-specification/gen/proto/task/v1"
)

func (g *GrpcClient) WithTaskManagementServiceClient(ctx context.Context) {
	log.Println("✅ Initializing TaskManagementService gRPC client")

	if g.config.TaskManagementService.Host == "" {
		log.Println("TaskManagementService is not configured")
		return
	}

	conn, err := newGrpcConnection(g.config.TaskManagementService)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to TaskManagementService: %v", err))
	}

	grpcClientList[g.config.TaskManagementService.Name] = conn

	g.MemberService = taskv1.NewMemberServiceClient(conn)

	log.Println("✅ TaskManagementService gRPC client initialized")
}
