package responses

type TestNotificationResponse struct {
	Message string `json:"message"`
}

type PaginationResponse struct {
	Page      int `json:"page"`
	PageSize  int `json:"pageSize"`
	TotalPage int `json:"totalPage"`
	TotalItem int `json:"totalItem"`
}
