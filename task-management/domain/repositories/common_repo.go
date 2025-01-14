package repositories

type PaginationRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	SortBy   string `json:"sortBy"`
	Order    string `json:"order"`
}
