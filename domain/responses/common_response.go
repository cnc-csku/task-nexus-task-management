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

type SetupStatusResponse struct {
	IsSetupOwner     bool `json:"isSetupOwner"`
	IsSetupWorkspace bool `json:"isSetupWorkspace"`
}

type GeneratePutPresignedURLResponse struct {
	URL       string `json:"url"`
	ExpiredIn string `json:"expiredIn"`
	ExpiredAt string `json:"expiredAt"`
}
