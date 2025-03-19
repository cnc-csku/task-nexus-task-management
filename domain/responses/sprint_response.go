package responses

import "time"

type CreateSprintResponse struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
}

type EditSprintResponse struct {
	Message string `json:"message"`
}

type CompleteSprintResponse struct {
	Message string `json:"message"`
}

type DeleteSprintResponse struct {
	Message string `json:"message"`
}
