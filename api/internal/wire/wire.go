//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/cnc-csku/task-nexus/internal/infrastructure/api"
	"github.com/google/wire"
)

func InitializeApp() *api.EchoAPI {
	wire.Build(
		CtxSet,
		ConfigSet,
		InfraSet,
		HandlerSet,
		api.NewEchoAPI,
	)

	return &api.EchoAPI{}
}
