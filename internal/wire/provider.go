package wire

import (
	core_grpcclient "github.com/cnc-csku/task-nexus-go-lib/grpcclient"
	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/domain/services"
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/repositories/grpcclient"
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/repositories/mongo"
	"github.com/cnc-csku/task-nexus/task-management/internal/adapters/rest"
	"github.com/cnc-csku/task-nexus/task-management/internal/infrastructure/cache"
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
	cache.NewRedisClient,
)

var RepositorySet = wire.NewSet(
	mongo.NewMongoUserRepo,
	mongo.NewMongoProjectRepo,
	mongo.NewMongoProjectMemberRepo,
	mongo.NewMongoWorkspaceRepo,
	mongo.NewMongoWorkspaceMemberRepo,
	mongo.NewMongoInvitationRepo,
	mongo.NewMongoGlobalSettingRepo,
	mongo.NewMongoSprintRepo,
	mongo.NewMongoTaskRepo,
	mongo.NewMongoTaskCommentRepo,
)

var ServiceSet = wire.NewSet(
	services.NewCommonService,
	services.NewUserService,
	services.NewProjectService,
	services.NewProjectMemberService,
	services.NewInvitationService,
	services.NewWorkspaceService,
	services.NewSprintService,
	services.NewTaskService,
	services.NewTaskCommentService,
)

var RestHandlerSet = wire.NewSet(
	rest.NewHealthCheckHandler,
	rest.NewCommonHandler,
	rest.NewUserHandler,
	rest.NewProjectHandler,
	rest.NewProjectMemberHandler,
	rest.NewInvitationHandler,
	rest.NewWorkspaceHandler,
	rest.NewSprintHandler,
	rest.NewTaskHandler,
	rest.NewTaskCommentHandler,
)

var GrpcClientSet = wire.NewSet(
	config.ProvideGrpcClientConfig,
	core_grpcclient.NewGrpcClient,
	grpcclient.NewGrpcClient,
)

var MiddlewareSet = wire.NewSet(
	middlewares.NewAdminJWTMiddleware,
)
