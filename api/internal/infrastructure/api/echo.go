package api

import (
	"context"
	"log"

	"github.com/cnc-csku/task-nexus/config"
	"github.com/cnc-csku/task-nexus/docs"
	"github.com/cnc-csku/task-nexus/internal/infrastructure/router"
	"github.com/labstack/echo/v4"
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

func (a *EchoAPI) Start() error {
	docs.SwaggerInfo.Title = "Task Nexus API"

	e := echo.New()

	a.router.RegisterAPIRouter(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	err := e.Start(":" + a.config.PORT)
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
