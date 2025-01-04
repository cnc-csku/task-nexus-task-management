package grpcclient

import (
	"fmt"

	notificationv1 "github.com/cnc-csku/task-nexus/api-specification/gen/proto/notification/v1"
	taskv1 "github.com/cnc-csku/task-nexus/api-specification/gen/proto/task/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	config GrpcClientConfig

	// task-management
	MemberService taskv1.MemberServiceClient

	// notification
	NotificationService notificationv1.NotificationServiceClient
}

var grpcClientList = make(map[string]*grpc.ClientConn)

func NewGrpcClient(
	config GrpcClientConfig,
) *GrpcClient {
	return &GrpcClient{
		config: config,
	}
}

func newGrpcConnection(cfg grpcClientConfig) (*grpc.ClientConn, error) {
	opts := initOptions(cfg)
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func initOptions(cfg grpcClientConfig) []grpc.DialOption {
	maxSendSize := cfg.MaxSendMsgSize * 1024 * 1024
	maxRecvSize := cfg.MaxRecvMsgSize * 1024 * 1024

	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(maxSendSize),
			grpc.MaxCallRecvMsgSize(maxRecvSize),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	return opts
}

func CloseAllGrpcConnections() {
	for _, conn := range grpcClientList {
		conn.Close()
	}
}
