package api

import (
	"context"

	taskv1 "github.com/cnc-csku/task-nexus/api-specification/gen/proto/task/v1"
	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/services"
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/grpcserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	Server *grpc.Server
	Config *config.Config
}

func NewGrpcServer(
	ctx context.Context,
	config *config.Config,
	memberService services.MemberService,
) *GrpcServer {
	opts := initOptions(config)

	server := grpc.NewServer(opts...)

	// Register health check
	healthServer := health.NewServer()
	healthServer.SetServingStatus(config.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	if config.GrpcUseReflection {
		reflection.Register(server)
	}

	taskv1.RegisterMemberServiceServer(server, &grpcserver.MemberHandler{MemberService: memberService})

	return &GrpcServer{
		Server: server,
		Config: config,
	}
}

func initOptions(
	config *config.Config,
) []grpc.ServerOption {
	const mbToBytes = 1024 * 1024
	maxSendSize := config.GrpcMaxSendSize * mbToBytes
	maxRecvSize := config.GrpcMaxRecvSize * mbToBytes

	opts := []grpc.ServerOption{
		grpc.MaxSendMsgSize(maxSendSize),
		grpc.MaxRecvMsgSize(maxRecvSize),
	}

	return opts
}
