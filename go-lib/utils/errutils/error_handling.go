package errutils

import "github.com/labstack/echo/v4"

type Error struct {
	Status  ErrorStatus
	Message string
}

func NewError(err error, errStatus ErrorStatus) *Error {
	return &Error{
		Status:  errStatus,
		Message: err.Error(),
	}
}

type RestErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (e *Error) ToEchoError() error {
	return echo.NewHTTPError(e.Status.StatusCode(), RestErrorResponse{
		Status:  e.Status.String(),
		Message: e.Message,
	})
}