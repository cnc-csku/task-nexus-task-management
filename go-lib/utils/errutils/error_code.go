package errutils

import "net/http"

type ErrorStatus string

const (
	BadRequest           ErrorStatus = "BAD_REQUEST"
	InternalError        ErrorStatus = "INTERNAL_ERROR"
	Unauthorized         ErrorStatus = "UNAUTHORIZED"
	Forbidden            ErrorStatus = "FORBIDDEN"
	NotFound             ErrorStatus = "NOT_FOUND"
	Conflict             ErrorStatus = "CONFLICT"
	UnsupportedMediaType ErrorStatus = "UNSUPPORTED_MEDIA_TYPE"
	UnprocessableEntity  ErrorStatus = "UNPROCESSABLE_ENTITY"
	TooManyRequests      ErrorStatus = "TOO_MANY_REQUESTS"
	InternalServerError  ErrorStatus = "INTERNAL_SERVER_ERROR"
)

func (e ErrorStatus) String() string {
	return string(e)
}

func (e ErrorStatus) StatusCode() int {
	statusMap := map[ErrorStatus]int{
		BadRequest:           http.StatusBadRequest,
		InternalError:        http.StatusInternalServerError,
		Unauthorized:         http.StatusUnauthorized,
		Forbidden:            http.StatusForbidden,
		NotFound:             http.StatusNotFound,
		Conflict:             http.StatusConflict,
		UnsupportedMediaType: http.StatusUnsupportedMediaType,
		UnprocessableEntity:  http.StatusUnprocessableEntity,
		TooManyRequests:      http.StatusTooManyRequests,
		InternalServerError:  http.StatusInternalServerError,
	}

	if code, exists := statusMap[e]; exists {
		return code
	}
	return http.StatusInternalServerError // Default status code
}