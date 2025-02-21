package grpcclient

import (
	"context"

	core_grpcclient "github.com/cnc-csku/task-nexus-go-lib/grpcclient"
)

type GrpcClient struct {
	Grpcclient *core_grpcclient.GrpcClient
}

func NewGrpcClient(
	ctx context.Context,
	grpcclient *core_grpcclient.GrpcClient,
) *GrpcClient {
	grpcclient.WithNotificationServiceClient(ctx)
	return &GrpcClient{
		Grpcclient: grpcclient,
	}
}
