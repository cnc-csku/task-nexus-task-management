package wire

import (
	core_grpcclient "github.com/cnc-csku/task-nexus/go-lib/grpcclient"
	"github.com/cnc-csku/task-nexus/notification/config"
	"github.com/cnc-csku/task-nexus/notification/domain/services"
	"github.com/cnc-csku/task-nexus/notification/internal/adapters/repositories/grpcclient"
	"github.com/google/wire"
)

var CtxSet = wire.NewSet(
	NewCtx,
)

var ConfigSet = wire.NewSet(
	config.NewConfig,
)

var ServiceSet = wire.NewSet(
	services.NewNotificationService,
)

var GrpcClientSet = wire.NewSet(
	config.ProvideGrpcClientConfig,
	core_grpcclient.NewGrpcClient,
	grpcclient.NewGrpcClient,
)
