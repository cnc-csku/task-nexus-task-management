package services

import (
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/repositories/grpcclient"
)

type CommonService interface{}

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
