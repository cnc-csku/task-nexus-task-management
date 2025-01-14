package api

import (
	"context"
	"log"

	"github.com/cnc-csku/task-nexus/go-lib/jsonvalidator"
	"github.com/cnc-csku/task-nexus/go-lib/logging"
	"github.com/cnc-csku/task-nexus/go-lib/utils/errutils"
	"github.com/cnc-csku/task-nexus/task-management/config"
	"github.com/cnc-csku/task-nexus/task-management/docs"
	"github.com/cnc-csku/task-nexus/task-management/internal/infrastructure/router"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
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

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: a.config.AllowOrigins,
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowMethods: []string{
			echo.GET,
			echo.PUT,
			echo.PATCH,
			echo.POST,
			echo.DELETE,
		},
	}))

	// Set up JSON validator
	e.Validator = jsonvalidator.NewValidator()

	// Custom error handler
	e.HTTPErrorHandler = errutils.CustomHTTPErrorHandler

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
