package wire

import (
	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/services"
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/repositories/mongo"
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/rest"
	"github.com/cnc-csku/task-nexus/task-management/internal/infrastructure/database"
	"github.com/cnc-csku/task-nexus/task-management/internal/infrastructure/router"
	"github.com/google/wire"
)

var CtxSet = wire.NewSet(
	NewCtx,
)

var ConfigSet = wire.NewSet(
	config.NewConfig,
)

var InfraSet = wire.NewSet(
	database.NewMongoClient,
	router.NewRouter,
)

var RepositorySet = wire.NewSet(
	mongo.NewMemberRepository,
)

var ServiceSet = wire.NewSet(
	services.NewMemberUseCase,
)

var RestHandlerSet = wire.NewSet(
	rest.NewHealthCheckHandler,
)
