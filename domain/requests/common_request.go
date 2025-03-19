package requests

type TestNotificationRequest struct{}

type PaginationRequest struct {
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"pageSize" query:"pageSize"`
	SortBy   string `json:"sortBy" query:"sortBy"`
	Order    string `json:"order" query:"order" validate:"oneof=ASC DESC asc desc ''"`
}

type GeneratePutPresignedURLRequest struct {
	FileName     string `json:"fileName" validate:"required"`
	FileCategory string `json:"fileCategory" validate:"required"`
}
