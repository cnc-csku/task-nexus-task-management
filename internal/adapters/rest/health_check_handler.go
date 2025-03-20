package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type healthCheckHandler struct{}

type HealthCheckHandler interface {
	HealthCheck(c echo.Context) error
}

func NewHealthCheckHandler() HealthCheckHandler {
	return &healthCheckHandler{}
}

// HealthCheck godoc
//
//	@Summary		Health Check
//	@Description	Check the health of the service
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object} map[string]string{message=string}
//	@Router			/api/health [get]
func (h *healthCheckHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, &echo.Map{
		"message": "OK",
	})
}
