package wire

import (
	core_grpcclient "github.com/cnc-csku/task-nexus/go-lib/grpcclient"
	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/services"
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/repositories/grpcclient"
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/repositories/mongo"
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/rest"
	"github.com/cnc-csku/task-nexus/task-management/internal/infrastructure/database"
	"github.com/cnc-csku/task-nexus/task-management/internal/infrastructure/llm"
	"github.com/cnc-csku/task-nexus/task-management/internal/infrastructure/router"
	"github.com/cnc-csku/task-nexus/task-management/middlewares"
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
	llm.NewOllamaClient,
)

var RepositorySet = wire.NewSet(
	mongo.NewMongoUserRepo,
	mongo.NewMongoProjectRepo,
)

var ServiceSet = wire.NewSet(
	services.NewCommonService,
	services.NewUserService,
	services.NewProjectService,
)

var RestHandlerSet = wire.NewSet(
	rest.NewHealthCheckHandler,
	rest.NewCommonHandler,
	rest.NewUserHandler,
	rest.NewProjectHandler,
)

var GrpcClientSet = wire.NewSet(
	config.ProvideGrpcClientConfig,
	core_grpcclient.NewGrpcClient,
	grpcclient.NewGrpcClient,
)

var MiddlewareSet = wire.NewSet(
	middlewares.NewAdminJWTMiddleware,
)
