package responses

import "time"

type UserResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	FullName    string    `json:"fullName"`
	DisplayName string    `json:"displayName"`
	ProfileUrl  string    `json:"profileUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserWithTokenResponse struct {
	UserResponse
	Token         string    `json:"token"`
	TokenExpireAt time.Time `json:"tokenExpireAt"`
}

type ListUserResponse struct {
	Users              []UserResponse     `json:"users"`
	PaginationResponse PaginationResponse `json:"pagination"`
}
