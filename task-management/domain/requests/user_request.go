package requests

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	FullName    string `json:"fullName" validate:"required"`
	DisplayName string `json:"displayName" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SearchUserParams struct {
	Keyword           string             `query:"keyword"`
	PaginationRequest *PaginationRequest `json:"pagination"`
}
