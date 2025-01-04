//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/cnc-csku/task-nexus/notification/internal/infrastructure/api"
	"github.com/google/wire"
)

func InitializeGrpcServer() *api.GrpcServer {
	wire.Build(
		CtxSet,
		ConfigSet,
		GrpcClientSet,
		ServiceSet,
		api.NewGrpcServer,
	)

	return &api.GrpcServer{}
}
