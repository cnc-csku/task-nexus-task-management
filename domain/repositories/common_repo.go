package repositories

type PaginationRequest struct {
	Page     int
	PageSize int
	SortBy   string
	Order    string
}
