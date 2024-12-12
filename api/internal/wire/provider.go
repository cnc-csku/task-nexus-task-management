package wire

import (
	"github.com/cnc-csku/task-nexus/config"
	"github.com/cnc-csku/task-nexus/internal/adapters/rest"
	"github.com/cnc-csku/task-nexus/internal/infrastructure/database"
	"github.com/cnc-csku/task-nexus/internal/infrastructure/router"
	"github.com/google/wire"
)

var ServiceSet = wire.NewSet()

var RepositorySet = wire.NewSet()

var HandlerSet = wire.NewSet(
	rest.NewHealthCheckHandler,
)

var InfraSet = wire.NewSet(
	database.NewMongoClient,
	router.NewRouter,
)

var CtxSet = wire.NewSet(
	NewCtx,
)

var ConfigSet = wire.NewSet(
	config.NewConfig,
)
