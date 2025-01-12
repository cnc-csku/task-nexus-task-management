package api

import (
	"context"
	"log"

	"github.com/cnc-csku/task-nexus/go-lib/logging"
	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/docs"
	"github.com/cnc-csku/task-nexus/task-management/internal/infrastructure/router"
	"github.com/cnc-csku/task-nexus/go-lib/jsonvalidator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EchoAPI struct {
	echo        *echo.Echo
	ctx         context.Context
	config      *config.Config
	mongoClient *mongo.Client
	router      *router.Router
}

func NewEchoAPI(
	ctx context.Context,
	config *config.Config,
	mongoClient *mongo.Client,
	router *router.Router,
) *EchoAPI {
	return &EchoAPI{
		echo:        echo.New(),
		ctx:         ctx,
		config:      config,
		mongoClient: mongoClient,
		router:      router,
	}
}

func (a *EchoAPI) Start(logger *logrus.Logger) error {
	docs.SwaggerInfo.Title = "Task Nexus API"

	e := echo.New()
	e.Use(logging.EchoLoggingMiddleware(logger))

	// Set up JSON validator
	e.Validator = jsonvalidator.NewValidator()

	a.router.RegisterAPIRouter(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	err := e.Start(":" + a.config.RestServer.Port)
	if err != nil {

		return err
	}

	defer func() {
		if err := a.mongoClient.Disconnect(a.ctx); err != nil {
			log.Printf("❌ Error disconnecting from MongoDB: %v\n", err)
		}
	}()

	return nil
}
