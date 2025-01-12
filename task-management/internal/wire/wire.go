//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/cnc-csku/task-nexus/task-management/internal/infrastructure/api"
	"github.com/google/wire"
)

func InitializeApp() *api.EchoAPI {
	wire.Build(
		CtxSet,
		ConfigSet,
		GrpcClientSet,
		InfraSet,
		RepositorySet,
		ServiceSet,
		RestHandlerSet,
		MiddlewareSet,
		api.NewEchoAPI,
	)

	return &api.EchoAPI{}
}

// func InitializeGrpcServer() *api.GrpcServer {
// 	wire.Build(
// 		CtxSet,
// 		ConfigSet,
// 		GrpcClientSet,
// 		InfraSet,
// 		RepositorySet,
// 		ServiceSet,
// 		api.NewGrpcServer,
// 	)

// 	return &api.GrpcServer{}
// }
